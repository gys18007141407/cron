package common

type Task struct {
	TaskName 				string		`json:"task_name"`
	TaskCommand 			string 		`json:"task_command"`
	TaskCronExpr			string 		`json:"task_cron_expr"`
}