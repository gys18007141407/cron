package main

import (
	"cron_master/src/config"
	"cron_master/src/master/logManager"
	"cron_master/src/master/router"
	"cron_master/src/master/taskManager"
	"cron_master/src/master/workerManager"
	"flag"
	"github.com/sirupsen/logrus"
	"runtime"
	"strconv"
)

// go本身是多线程的，我们在程序中创建的是协程。为了让go效率最大化，设置线程数量等于内核数量
func init()  {
	var(
		cpus 		int
	)
	cpus = runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	logrus.Infoln("set threads num =", cpus)
}

var(
	configFile  			string
)

func main()  {
	var (
		err error
	)

	// 解析命令行参数
	flag.StringVar(&configFile, "config", "/home/gys/go/src/cron_master/src/config/config.ini", "指定配置文件路径")
	flag.Parse()

	// 读取并初始化配置
	if err = config.InitConfig(configFile); err != nil{
		logrus.Fatalln(err)
		return
	}
	// 打印配置信息
	logrus.Infoln(*config.Cfg)

	// 初始化任务管理器
	if err = taskManager.InitTaskManager(); err != nil{
		logrus.Fatalln(err)
		return
	}

	// 初始化日志管理器
	if err = logManager.InitLogManager(); err != nil{
		logrus.Fatalln(err)
		return
	}

	// 初始化worker集群管理器(服务发现)
	if err = workerManager.InitWorkerManager(); err != nil{
		logrus.Fatalln(err)
		return
	}

	// 启动路由
	if err = router.Router.Run(config.Cfg.ListenIP + ":" + strconv.Itoa(config.Cfg.ListenPort)); err != nil{
		logrus.Fatalln(err)
		return
	}

	logrus.Infoln("退出中.....")
}
