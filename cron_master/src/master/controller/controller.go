package controller

import (
	"cron_master/src/common"
	"cron_master/src/master/logManager"
	"cron_master/src/master/taskManager"
	"cron_master/src/master/workerManager"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// POST 保存任务到etcd中
func SaveTask(c *gin.Context)  {
	var(
		err 			error
		task 			common.Task
		oldTask 		*common.Task
	)
	if err = c.BindJSON(&task); err != nil{
		c.JSON(http.StatusCreated, gin.H{
			"errno":1,
			"message":err.Error(),
			"data":nil,
		})
		return
	}

	// 通过任务管理器添加任务
	if oldTask, err = taskManager.TM.SaveTask(&task); err != nil{
		c.JSON(http.StatusAccepted, gin.H{
			"errno":1,
			"message":err.Error(),
			"data":nil,
		})
	}else{
		c.JSON(http.StatusOK, gin.H{
			"errno":0,
			"message":"success",
			"data":oldTask,
		})
	}
}

// DELETE 从etcd中删除任务
func RemoveTask(c *gin.Context)  {
	var (
		ok					bool
		err 				error
		taskName 			string
		oldTask 			*common.Task
	)
	if taskName, ok = c.GetQuery("task_name"); !ok{
		c.JSON(http.StatusCreated, gin.H{
			"errno":1,
			"message":"缺少query字段:task_name",
			"data":nil,
		})
		return
	}

	if oldTask, err = taskManager.TM.RemoveTask(taskName); err != nil{
		c.JSON(http.StatusAccepted, gin.H{
			"errno":1,
			"message":"success",
			"data":oldTask,
		})
	}else{
		c.JSON(http.StatusOK, gin.H{
			"errno":0,
			"message":"success",
			"data":oldTask,
		})
	}
}


// GET 获取etcd当前有哪些任务
func GetTasks(c *gin.Context)  {
	var (
		err 			error
		taskList 		[]*common.Task
	)

	if taskList, err = taskManager.TM.GetTaskList(); err != nil{
		c.JSON(http.StatusAccepted, gin.H{
			"errno":1,
			"message":err.Error(),
			"data":nil,
		})
	}else{
		c.JSON(http.StatusOK, gin.H{
			"errno":0,
			"message":"success",
			"data":taskList,
		})
	}
}

// POST 强制杀死任务进程
func KillTask(c *gin.Context)  {
	var (
		ok					bool
		err 				error
		taskName 			string
	)
	if taskName, ok = c.GetQuery("task_name"); !ok{
		c.JSON(http.StatusCreated, gin.H{
			"errno":1,
			"message":"缺少query字段:task_name",
		})
		return
	}

	if err = taskManager.TM.KillTask(taskName); err != nil{
		c.JSON(http.StatusAccepted, gin.H{
			"errno":1,
			"message":err.Error(),
		})
	}else{
		c.JSON(http.StatusOK, gin.H{
			"errno":0,
			"message":"success",
		})
	}
}

// GET 获取任务执行日志
func QueryTaskLog(c *gin.Context) {
	var(
		ok 				bool
		err 			error
		taskName		string
		skipStr			string
		skip 			int
		limitStr		string
		limit 			int
		taskLogList		[]*common.TaskLog

	)
	if taskName, ok = c.GetQuery("task_name"); !ok{
		c.JSON(http.StatusCreated, gin.H{
			"errno":1,
			"message":"缺少query字段:task_name",
			"data":nil,
		})
		return
	}
	if skipStr, ok = c.GetQuery("skip"); !ok{
		c.JSON(http.StatusCreated, gin.H{
			"errno":1,
			"message":"缺少query字段:skip",
			"data":nil,
		})
		return
	}
	if limitStr, ok = c.GetQuery("limit"); !ok{
		c.JSON(http.StatusCreated, gin.H{
			"errno":1,
			"message":"缺少query字段:limit",
			"data":nil,
		})
		return
	}

	if skip, err = strconv.Atoi(skipStr); err != nil{
		c.JSON(http.StatusCreated, gin.H{
			"errno":1,
			"message": err.Error(),
			"data":nil,
		})
		return
	}
	if limit, err = strconv.Atoi(limitStr); err != nil{
		c.JSON(http.StatusCreated, gin.H{
			"errno":1,
			"message":err.Error(),
			"data":nil,
		})
		return
	}


	if taskLogList, err = logManager.LM.QueryTaskLog(taskName, int64(skip), int64(limit)); err != nil{
		c.JSON(http.StatusAccepted, gin.H{
			"errno":1,
			"message":err.Error(),
			"data":nil,
		})
	}else{
		c.JSON(http.StatusOK, gin.H{
			"errno":0,
			"message":"success",
			"data":taskLogList,
		})
	}
}

// GET 获取健康的worker节点
func GetWorkers(c *gin.Context)  {
	var (
		err 			error
		workers			[]string
	)

	if workers, err = workerManager.WM.GetWorkers(); err != nil{
		c.JSON(http.StatusAccepted, gin.H{
			"errno":1,
			"message":err.Error(),
			"data":nil,
		})
	}else{
		c.JSON(http.StatusOK, gin.H{
			"errno":0,
			"message":"success",
			"data":workers,
		})
	}
}