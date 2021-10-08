package scheduler

import (
	"context"
	"cron_worker/src/common"
	"cron_worker/src/worker/executor"
	"cron_worker/src/worker/logger"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"time"
)

// 任务调度器
type Scheduler struct {
	// 任务事件队列
	EventChan 			chan *common.TaskEvent
	// 调度计划表
	SchedulePlan 		map[string]*common.TaskSchedule
	// 任务执行状态表
	ExecStatus			map[string]*common.TaskExecStatus
	// 任务执行结果队列
	ExecResultChan		chan *common.TaskExecResult
}

// 调度器的事件循环:监听调度器管道
func (This *Scheduler) loop ()  {
	var(
		err 				error
		taskEvent			*common.TaskEvent
		taskExecResult		*common.TaskExecResult
		nearestExpired		time.Duration

	)

	// 初始化(当前没有任务)
	nearestExpired = This.UpdateTaskSchedule()

	// 处理到来的调度事件
	for{
		select {
		case taskEvent = <-This.EventChan:
			if err = This.solveTaskEvent(taskEvent); err != nil{
				logrus.Infoln(err)
			}
			break
		case taskExecResult = <-This.ExecResultChan:
			if err = This.solveTaskExecResult(taskExecResult); err != nil{
				logrus.Infoln(err)
			}
			break
		case <-time.After(nearestExpired):
			// logrus.Infoln("精准睡眠，定时器到期")
			break
		}
		nearestExpired = This.UpdateTaskSchedule()
	}
}

// 调度器处理任务事件
func (This *Scheduler) solveTaskEvent(taskEvent *common.TaskEvent) (err error)  {
	var(
		taskSchedule 	*common.TaskSchedule
		taskExecStatus	*common.TaskExecStatus
		ok				bool
	)

	switch taskEvent.CurEvent {
	case common.EventSave:
		// 更新调度计划表中的该任务
		if taskSchedule, err = This.NewTaskSchedule(taskEvent.CurTask); err != nil{
			return
		}else{
			This.SchedulePlan[taskEvent.CurTask.TaskName] = taskSchedule
		}
		break
	case common.EventDelete:
		// 从调度计划表中删除该任务
		if taskSchedule, ok = This.SchedulePlan[taskEvent.CurTask.TaskName]; !ok{
			return common.ERROR_SCHEDULEPLAN
		}else{
			delete(This.SchedulePlan, taskEvent.CurTask.TaskName)
		}
		break
	case common.EventKill:
		// 强杀该任务(取消Command执行, CancelFunc())
		if taskExecStatus, ok = This.ExecStatus[taskEvent.CurTask.TaskName]; !ok{
			return common.ERROR_KILLTASK
		}else{
			taskExecStatus.DoCancelFunc()
			// delete(This.ExecStatus, taskEvent.CurTask.TaskName)
		}
		break

	default:
		return common.ERROR_TASKEVENT
	}
	return nil
}

// 调度器处理任务完成信息[写任务日志]
func (This *Scheduler) solveTaskExecResult(taskExecResult *common.TaskExecResult) (err error) {
	var (
		ok			bool
	)
	logrus.Infoln(taskExecResult.CurTaskExecStatus.CurTask.TaskName, "output=", string(taskExecResult.CurTaskOutput), "err=", taskExecResult.CurTaskError)

	if _, ok = This.ExecStatus[taskExecResult.CurTaskExecStatus.CurTask.TaskName]; !ok{
		return common.ERROR_DELETE_EXECSTATUS
	}
	delete(This.ExecStatus, taskExecResult.CurTaskExecStatus.CurTask.TaskName)
	// 写日志
	if taskExecResult.CurTaskError != common.ERROR_LOCK_REQUIRED {
		logger.Logger.PushTaskLog(This.NewTaskLog(taskExecResult))
	}
	return nil
}

// 重新计算所有任务的状态
func (This *Scheduler) UpdateTaskSchedule() (nearestExpired time.Duration) {
	var(
		now 					time.Time
		taskSchedule 			*common.TaskSchedule
		nearest 				*time.Time
	)
	now = time.Now()

	// 没有任务
	if len(This.SchedulePlan) == 0{
		return time.Second
	}

	// 遍历所有任务计划
	for _, taskSchedule = range This.SchedulePlan{
		// 任务到期
		if taskSchedule.NextTime.Before(now) || taskSchedule.NextTime.Equal(now){
			// 尝试执行任务(当执行时间大于调度间隔时，任务可能还在执行)
			This.ExecTask(taskSchedule)
			taskSchedule.NextTime = taskSchedule.CronSched.Next(now)
		}

		if nearest == nil || taskSchedule.NextTime.Before(*nearest){
			nearest = &taskSchedule.NextTime
		}
	}

	// 精准睡眠
	if nearest.After(now){
		nearestExpired = nearest.Sub(now)
	}
	return
}

// 执行任务[调度 != 执行]
func (This *Scheduler) ExecTask (taskSchedule *common.TaskSchedule)  {
	var(
		ok						bool						// isExecuting
		taskExecStatus 			*common.TaskExecStatus
	)

	if taskExecStatus, ok = This.ExecStatus[taskSchedule.CurTask.TaskName]; ok{
		logrus.Infoln("任务正在执行，上一次执行还未完成，因此跳过本次执行....")
		return
	}

	// 新建任务执行状态并保存
	taskExecStatus = This.NewTaskExecStatus(taskSchedule)
	This.ExecStatus[taskSchedule.CurTask.TaskName] = taskExecStatus

	// 将任务提交给执行器
	// executor.Exec.ExecTask(taskExecStatus)			// import cycle
	executor.Exec.ExecTask(taskExecStatus, This.PushTaskExecResult)
}


// 给调度器推送任务事件
func (This *Scheduler) PushTaskEvent (taskEvent *common.TaskEvent)  {
	select {
	case This.EventChan <- taskEvent:
		break
	default:
		logrus.Infoln("任务事件队列满,已经忽略该任务事件")
		break
	}
}

// 给调度器推送任务执行结果
func (This *Scheduler) PushTaskExecResult (taskExecResult *common.TaskExecResult)  {
	This.ExecResultChan <- taskExecResult
}

// 根据Task任务创建新的TaskSchedule任务调度计划
func (This *Scheduler) NewTaskSchedule (task *common.Task) (taskSchedule *common.TaskSchedule, err error)  {
	var(
		cronSchedule 		cron.Schedule
	)

	// 解析cron表达式
	if cronSchedule, err = cron.Parse(task.TaskCronExpr); err != nil{
		return
	}

	taskSchedule = &common.TaskSchedule{
		CurTask:   task,
		CronSched: cronSchedule,
		NextTime:  cronSchedule.Next(time.Now()),
	}
	return
}

// 根据TaskSchedule任务调度计划创建新的TaskExecStatus任务执行状态信息
func (This *Scheduler) NewTaskExecStatus (taskSchedule *common.TaskSchedule) (taskExecStatus *common.TaskExecStatus) {
	var(
		ctx 			context.Context
		cancelFunc 		context.CancelFunc
	)
	// TODO: Task增设超时选项，这里ctx继承withTimeout上下文，不仅可以手动强杀任务,超时也自动取消执行任务。
	ctx, cancelFunc = context.WithCancel(context.TODO())
	taskExecStatus = &common.TaskExecStatus{
		CurTask:          taskSchedule.CurTask,
		ScheduleTime:     taskSchedule.NextTime,
		RealScheduleTime: time.Now(),
		ExecTime:         time.Time{},
		FinishTime:       time.Time{},
		CancelCtx:        ctx,
		DoCancelFunc:     cancelFunc,
	}
	return
}


// 创建新的任务日志
func (This *Scheduler) NewTaskLog(taskExecResult *common.TaskExecResult) (taskLog *common.TaskLog) {
	taskLog = &common.TaskLog{
		TaskName:         taskExecResult.CurTaskExecStatus.CurTask.TaskName,
		TaskCommand:      taskExecResult.CurTaskExecStatus.CurTask.TaskCommand,
		TaskOutput:       string(taskExecResult.CurTaskOutput),
		ScheduleTime:     taskExecResult.CurTaskExecStatus.ScheduleTime.UnixNano() / 1000 / 1000,
		RealScheduleTime: taskExecResult.CurTaskExecStatus.RealScheduleTime.UnixNano() / 1000 / 1000,
		ExecTime:         taskExecResult.CurTaskExecStatus.ExecTime.UnixNano() / 1000 / 1000,
		FinishTime:       taskExecResult.CurTaskExecStatus.FinishTime.UnixNano() / 1000 / 1000,
	}
	if taskExecResult.CurTaskError != nil{
		taskLog.TaskError = taskExecResult.CurTaskError.Error()
	}
	return
}

// 任务调度器单例
var (
	Sched		*Scheduler
)

func init()  {
	Sched = &Scheduler{
		EventChan: make(chan *common.TaskEvent, 512),
		SchedulePlan: make(map[string]*common.TaskSchedule, 512),
		ExecStatus: make(map[string]*common.TaskExecStatus, 512),
		ExecResultChan: make(chan *common.TaskExecResult, 512),
	}

	// 启动任务调度器
	go Sched.loop()
}