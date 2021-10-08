package common

import (
	"github.com/robfig/cron"
	"time"
)

type TaskSchedule struct{
	CurTask				*Task					// 调度任务信息
	CronSched 			cron.Schedule			// 该任务的cron
	NextTime 			time.Time				// 下一次调度的时间
}