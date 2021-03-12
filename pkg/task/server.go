package task

/*
  @Author : lanyulei
*/

import (
	"ferry/pkg/logger"
	"ferry/pkg/task/worker"
)

func Start() {
	// 1. 启动服务，连接redis
	worker.StartServer()

	// 2. 启动异步调度
	taskWorker := worker.NewAsyncTaskWorker(10)
	err := taskWorker.Launch()
	if err != nil {
		logger.Errorf("启动machinery失败，%v", err.Error())
	}
}
