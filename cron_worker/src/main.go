package main

import (
	"cron_worker/src/config"
	"cron_worker/src/worker/logger"
	"cron_worker/src/worker/register"
	"cron_worker/src/worker/taskManager"
	"flag"
	"github.com/sirupsen/logrus"
	"runtime"
	"time"
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
	flag.StringVar(&configFile, "config", "/home/gys/go/src/cron_worker/src/config/config.ini", "指定配置文件路径")
	flag.Parse()

	// 读取并初始化配置
	if err = config.InitConfig(configFile); err != nil{
		logrus.Fatalln(err)
		return
	}

	// 初始化etcd客户端
	if err = taskManager.InitTaskManager(); err != nil{
		logrus.Fatalln(err)
		return
	}

	// 初始化日志记录器
	if err = logger.InitLogger(); err != nil{
		logrus.Fatalln(err)
		return
	}

	// 初始化服务注册器
	if err = register.InitWorkerRegister(); err != nil{
		logrus.Fatalln(err)
		return
	}

	// 打印配置信息
	logrus.Infoln(*config.Cfg)

	// 开始监听任务目录
	go func() {
		if err = taskManager.TM.WatchTasks(); err != nil{
			logrus.Fatalln(err)
			return
		}
	}()

	// 开始监听强杀目录
	go func() {
		if err = taskManager.TM.WatchKiller(); err != nil{
			logrus.Fatalln(err)
			return
		}
	}()

	for{
		select {
		case <-time.After(60*time.Second):
			break
		}
	}
}
