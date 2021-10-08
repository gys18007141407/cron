package executor

import (
	"cron_worker/src/common"
	"cron_worker/src/worker/lock"
	"os/exec"
	"time"
)

// 任务执行器
type Executor struct {

}

// 绑定的方法
// 执行任务
func (This *Executor) ExecTask (taskExecStatus *common.TaskExecStatus, callback func(result *common.TaskExecResult)())  {
	// 创建一个协程来执行该任务
	go func() {
		var (
			err						error
			output					[]byte
			taskExecResult			*common.TaskExecResult

			cmd						*exec.Cmd

			taskLock 				*lock.TaskLock
		)
		// 分布式锁
		taskLock = lock.NewLock(taskExecStatus.CurTask.TaskName)
		if err = taskLock.TryLock(); err != nil{
			goto CREATE_EXEC_RESULT
		}

		// 新建cmd
		cmd = exec.CommandContext(taskExecStatus.CancelCtx, "bash", "-c", taskExecStatus.CurTask.TaskCommand)

		taskExecStatus.ExecTime = time.Now()
		// 执行cmd
		output, err = cmd.CombinedOutput()
		taskExecStatus.FinishTime = time.Now()

CREATE_EXEC_RESULT:
		// 解锁
		taskLock.UnLock()

		// 执行结果信息
		taskExecResult = &common.TaskExecResult{
			CurTaskExecStatus: taskExecStatus,
			CurTaskOutput:     output,
			CurTaskError:      err,
		}

		// 将执行结果传给调度器(协程通信)
		// scheduler.Sched.PushTaskExecResult(taskExecResult)  // import cycle
		callback(taskExecResult)
	}()
}


// 任务执行器单例
var (
	Exec		*Executor
)

func init()  {
	Exec = &Executor{}
}