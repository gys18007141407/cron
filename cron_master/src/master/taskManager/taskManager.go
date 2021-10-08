package taskManager

import (
	"context"
	"cron_master/src/common"
	"cron_master/src/config"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"path"
)

// 一个etcd客户端，用来管理任务
type TaskManager struct {
	client 		*clientv3.Client
	kv 			clientv3.KV
	lease 		clientv3.Lease
}

// 该客户端绑定的方法
func (This *TaskManager) SaveTask(task *common.Task) (oldTask *common.Task, err error) {
	// 将该任务task保存到/cron/tasks/taskname

	var(
		taskKey			string
		taskValue		[]byte

		temp 			common.Task

		Op 				clientv3.Op
		OpResp			clientv3.OpResponse
	)

	taskKey = path.Join(config.Cfg.BaseDir, task.TaskName)
	if taskValue, err = json.Marshal(task); err != nil{
		return
	}

	Op = clientv3.OpPut(taskKey, string(taskValue), clientv3.WithPrevKV())
	if OpResp, err = This.kv.Do(context.TODO(), Op); err != nil{
		return
	}

	// 如果是更新任务的话，反序列化原来的任务
	if OpResp.Put().PrevKv != nil{
		if err = json.Unmarshal(OpResp.Put().PrevKv.Value, &temp); err != nil{
			logrus.Infoln("SaveTask反序列化错误...已丢弃该错误:", err.Error())
			return nil, nil
		}
		oldTask = &temp
	}
	return
}

func (This *TaskManager) RemoveTask(taskName string) (oldTask *common.Task, err error) {
	// 从etcd中移除/cron/tasks/taskName
	var (
		taskKey			string
		Op				clientv3.Op
		OpResp			clientv3.OpResponse
		temp			common.Task
	)

	// 该任务在etcd中的key
	taskKey = path.Join(config.Cfg.BaseDir, taskName)

	Op = clientv3.OpDelete(taskKey, clientv3.WithPrevKV())
	if OpResp, err = This.kv.Do(context.TODO(), Op); err != nil{
		return
	}

	// 如果删除成功, 反序列化原来的任务
	if len(OpResp.Del().PrevKvs) != 0{
		if err = json.Unmarshal(OpResp.Del().PrevKvs[0].Value, &temp); err != nil{
			logrus.Infoln("RemoveTask反序列化错误...已丢弃该错误:", err.Error())
			return nil, nil
		}
		oldTask = &temp
	}
	return
}

func (This *TaskManager) GetTaskList () (taskList []*common.Task, err error) {
	// 查询/cron/tasks/目录下的所有key
	var (
		Op 				clientv3.Op
		OpResp			clientv3.OpResponse
		kvPair			*mvccpb.KeyValue
		temp			*common.Task
	)
	taskList = make([]*common.Task, 0)

	// 获取任务根目录下的所有任务
	Op = clientv3.OpGet(config.Cfg.BaseDir, clientv3.WithPrevKV(), clientv3.WithPrefix())
	if OpResp, err = This.kv.Do(context.TODO(), Op); err != nil{
		return
	}

	// 遍历Response, 并进行反序列化
	for _, kvPair = range OpResp.Get().Kvs{
		temp = &common.Task{}
		if err = json.Unmarshal(kvPair.Value, temp); err != nil {
			logrus.Infoln("GetTaskList反序列化错误...已丢弃该错误:", err.Error())
			temp.TaskName = "反序列化错误"
			temp.TaskCommand = "?"
			temp.TaskCronExpr = "?"
			taskList = append(taskList, temp)
		}else{
			taskList = append(taskList, temp)
		}
	}

	return taskList, nil
}

func (This * TaskManager) KillTask (taskName string) (err error){
	// 更新etcd中该任务，将key更新到killer目录下（worker监听该目录，杀死任务）
	var (
		taskKillerKey	string

		leaseGrantResp	*clientv3.LeaseGrantResponse
		leaseID			clientv3.LeaseID

		Op 				clientv3.Op
	)
	taskKillerKey = path.Join(config.Cfg.KillerDir, taskName)

	// 申请租约(该key在强杀目录中存在2s，然后被删除。先触发PUT事件，然后触发删除事件)
	if leaseGrantResp, err = TM.lease.Grant(context.TODO(), 2); err != nil{
		return
	}

	// 该租约的ID
	leaseID = leaseGrantResp.ID

	// 置killer标记：将其put到killer目录下，worker监听到后强杀该任务
	Op = clientv3.OpPut(taskKillerKey, "", clientv3.WithLease(leaseID))
	if _, err = This.kv.Do(context.TODO(), Op); err != nil{
		return
	}
	logrus.Infoln("杀死任务：", taskName)
	return
}


// 任务管理器单例
var(
	TM 				*TaskManager
)

// 初始化单例
func InitTaskManager() (err error) {
	if TM == nil {
		var (
			etcdConfig 		clientv3.Config
			client 			*clientv3.Client
		)
		etcdConfig = clientv3.Config{
			Endpoints:            config.Cfg.Endpoints,
			DialTimeout:          config.Cfg.DialTimeout,
		}
		
		// 建立连接
		if client, err = clientv3.New(etcdConfig); err != nil{
			return err
		}
		

		// 赋值单例
		TM = &TaskManager{
			client: client,
			kv:     clientv3.NewKV(client),
			lease:  clientv3.NewLease(client),
		}
	}
	return nil
}
