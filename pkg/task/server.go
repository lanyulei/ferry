package task

/*
  @Author : lanyulei
*/

import (
	"ferry/pkg/logger"
	"ferry/pkg/task/worker"
)

func Start() {
	// 启动异步任务框架
	taskWorker := worker.NewAsyncTaskWorker(0)
	err := taskWorker.Launch()
	if err != nil {
		logger.Errorf("启动machinery失败，%v", err.Error())
	}
}
