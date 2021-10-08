package common

type TaskLog struct {
	TaskName 					string		`bson:"task_name" json:"task_name"`							// 任务名称
	TaskCommand 				string 		`bson:"task_command" json:"task_command"`					// 任务指令
	TaskOutput					string		`bson:"task_output" json:"task_output"`						// 任务标准输出
	TaskError					string		`bson:"task_error" json:"task_error"`						// 任务错误输出
	ScheduleTime				int64		`bson:"schedule_time" json:"schedule_time"`					// 理论被调度的时间
	RealScheduleTime			int64		`bson:"real_schedule_time" json:"real_schedule_time"`		// 真正被调度的时间
	ExecTime					int64		`bson:"exec_time" json:"exec_time"`							// 执行时间
	FinishTime 					int64		`bson:"finish_time" json:"finish_time"`						// 完成时间
}
