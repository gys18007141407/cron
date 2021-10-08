package router

import (
	"cron_master/src/master/controller"
	"cron_master/src/master/middleware"
	"github.com/gin-gonic/gin"
)

// 路由单例(gin引擎单例)
var(
	Router			*gin.Engine
)

func init() {
	if Router == nil {
		var (
			adminRouter *gin.RouterGroup
		)
		// 初始化路由(初始化gin引擎)
		Router = gin.Default()
		// 全局路由中间件
		Router.Use(middleware.Cors())

		// 路由匹配
		adminRouter = Router.Group("/api/v1")
		{
			adminRouter.POST("/save", controller.SaveTask)

			adminRouter.DELETE("/delete", controller.RemoveTask)

			adminRouter.GET("/list", controller.GetTasks)

			adminRouter.POST("/kill", controller.KillTask)

			adminRouter.GET("/log", controller.QueryTaskLog)

			adminRouter.GET("/worker", controller.GetWorkers)
		}
	}
}
