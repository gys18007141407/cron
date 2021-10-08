package logManager

import (
	"context"
	"cron_master/src/common"
	"cron_master/src/config"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogManager struct {
	mongoClient 			*mongo.Client
	mongoCollection			*mongo.Collection
}


func (This *LogManager) newTaskNameFilter(taskName string) bson.D {
	return bson.D{{"task_name", taskName}}
}

func (This *LogManager) QueryTaskLog (taskName string, skip int64, limit int64) (logList []*common.TaskLog ,err error) {
	var(
		cursor			*mongo.Cursor
		findOpt			*options.FindOptions
		taskLog 		*common.TaskLog
	)
	logList = make([]*common.TaskLog, 0)

	findOpt = &options.FindOptions{
		Limit:               &limit,
		Skip:                &skip,
		Sort:                bson.D{{"exec_time", -1}},
	}
	if cursor, err = This.mongoCollection.Find(context.TODO(), This.newTaskNameFilter(taskName), findOpt); err != nil{
		return
	}
	defer func() {
		if err = cursor.Close(context.TODO()); err != nil{
			return
		}
	}()

	for cursor.Next(context.TODO()){
		taskLog = &common.TaskLog{}
		// 反序列化
		if err = cursor.Decode(taskLog); err != nil{
			logrus.Infoln("一条日志反序列化失败，已忽略该条日志:", err)
		}else{
			logList = append(logList, taskLog)
		}
	}

	return
}


// 日志管理器单例
var (
	LM				*LogManager
)

func InitLogManager() (err error) {
	if LM == nil{
		var (
			client 			*mongo.Client
			ctx        		context.Context
			cancelFunc 		context.CancelFunc
		)

		ctx, cancelFunc = context.WithTimeout(context.TODO(), config.Cfg.ConnectTimeOut)
		defer cancelFunc()

		if client, err = mongo.Connect(ctx, options.Client().ApplyURI(config.Cfg.DatabaseURI)); err != nil {
			return err
		}

		// 赋值单例
		LM = &LogManager{
			mongoClient:     client,
			mongoCollection: client.Database(config.Cfg.DatabaseName).Collection(config.Cfg.Collection),
		}
	}
	return nil
}