package config

import (
	"github.com/Unknwon/goconfig"
	"strconv"
	"strings"
	"time"
)

// 加载的配置
type Config struct {
	// web
	ListenIP 			string
	ListenPort 			int

	// etcd
	Endpoints 			[]string
	DialTimeout 		time.Duration

	// task
	BaseDir 			string
	KillerDir			string

	// worker
	WorkersDir			string

	// database
	DatabaseURI			string
	ConnectTimeOut		time.Duration
	DatabaseName 		string
	Collection 			string
}

// 配置的单例
var(
	Cfg 	*Config
)

// 初始化配置
func InitConfig(configFile string) error {
	if Cfg == nil {
		var (
			err 		error
			cf 			*goconfig.ConfigFile
			config		Config
		)
		// 读取配置文件
		if cf, err = goconfig.LoadConfigFile(configFile); err != nil {
			return err
		}

		// 初始化

		if err = initWebConfig(cf, &config); err != nil{
			return err
		}

		if err = initEtcdConfig(cf, &config); err != nil{
			return err
		}

		if err = initTaskConfig(cf, &config); err != nil{
			return err
		}

		if err = initWorkerConfig(cf, &config); err != nil{
			return err
		}

		if err = initDatabaseConfig(cf, &config); err != nil{
			return err
		}

		Cfg = &config
	}
	return nil
}

// 初始化web配置
func initWebConfig(cf *goconfig.ConfigFile, config *Config) (err error) {
	var(
		port 		string
	)
	if config.ListenIP, err = cf.GetValue("web", "ip"); err != nil{
		return err
	}

	if port, err = cf.GetValue("web", "port"); err != nil{
		return err
	}

	if config.ListenPort, err = strconv.Atoi(port); err != nil{
		return err
	}

	return nil
}


// 初始etcd配置
func initEtcdConfig(cf *goconfig.ConfigFile, config *Config) (err error) {
	var(
		endpoints 			string
		dialTimeoutStr 		string
		dialTimeout 		int
	)

	if endpoints, err = cf.GetValue("etcd", "Endpoints"); err != nil{
		return err
	}

	if dialTimeoutStr, err = cf.GetValue("etcd", "DialTimeout"); err != nil{
		return err
	}

	if dialTimeout, err = strconv.Atoi(dialTimeoutStr); err != nil{
		return err
	}

	config.Endpoints = strings.Split(endpoints, ",")
	config.DialTimeout = time.Duration(dialTimeout)*time.Millisecond

	return nil
}

// 初始task配置
func initTaskConfig(cf *goconfig.ConfigFile, config *Config) (err error) {
	var(
		baseDir 			string
		killerDir 			string
	)

	if baseDir, err = cf.GetValue("task", "BaseDir"); err != nil{
		return err
	}
	if killerDir, err = cf.GetValue("task", "KillerDir"); err != nil{
		return err
	}

	config.BaseDir = baseDir
	config.KillerDir = killerDir

	return nil
}

// 初始worker配置
func initWorkerConfig(cf *goconfig.ConfigFile, config *Config) (err error) {
	var(
		workersDir			string
	)

	if workersDir, err = cf.GetValue("worker", "WorkersDir"); err != nil{
		return err
	}

	config.WorkersDir = workersDir

	return nil
}

// 初始mongodb配置
func initDatabaseConfig(cf *goconfig.ConfigFile, config *Config) (err error) {
	var(
		dbURI 					string
		connectTimeOutStr 		string
		connectTimeOut			int
		dbName					string
		collection 				string
	)

	if dbURI, err = cf.GetValue("database", "DatabaseURI"); err != nil{
		return err
	}
	if connectTimeOutStr, err = cf.GetValue("database", "ConnectTimeOut"); err != nil{
		return err
	}
	if dbName, err = cf.GetValue("database", "DatabaseName"); err != nil{
		return err
	}
	if collection, err = cf.GetValue("database", "Collection"); err != nil{
		return err
	}

	if connectTimeOut, err = strconv.Atoi(connectTimeOutStr); err != nil{
		return err
	}

	config.DatabaseURI = dbURI
	config.ConnectTimeOut = time.Duration(connectTimeOut)*time.Millisecond
	config.DatabaseName = dbName
	config.Collection = collection

	return nil
}
