package worker

import (
	"github.com/RichardKnop/machinery/v1"
	taskConfig "github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
)

var AsyncTaskCenter *machinery.Server

func init() {
	tc, err := NewTaskCenter()
	if err != nil {
		panic(err)
	}
	AsyncTaskCenter = tc
}

func NewTaskCenter() (*machinery.Server, error) {
	cnf := &taskConfig.Config{
		Broker:        "redis://127.0.0.1:6379",
		DefaultQueue:  "ServerTasksQueue",
		ResultBackend: "eager",
	}
	server, err := machinery.NewServer(cnf)
	if err != nil {
		return nil, err
	}
	initAsyncTaskMap()
	return server, server.RegisterTasks(asyncTaskMap)
}

func NewAsyncTaskWorker(concurrency int) *machinery.Worker {
	consumerTag := "TaskWorker"
	worker := AsyncTaskCenter.NewWorker(consumerTag, concurrency)
	errorHandler := func(err error) {
		log.ERROR.Println("执行失败: ", err)
	}
	preTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("开始执行: ", signature.Name)
	}
	postTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("执行结束: ", signature.Name)
	}
	worker.SetPostTaskHandler(postTaskHandler)
	worker.SetErrorHandler(errorHandler)
	worker.SetPreTaskHandler(preTaskHandler)
	return worker
}
