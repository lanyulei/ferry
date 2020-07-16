package task

/*
  @Author : lanyulei
*/

import (
	"ferry/pkg/task/worker"

	"github.com/RichardKnop/machinery/v1/log"
)

func Start() {
	// 启动异步任务框架
	taskWorker := worker.NewAsyncTaskWorker(0)
	err := taskWorker.Launch()
	if err != nil {
		log.ERROR.Println("启动machinery失败，%v", err.Error())
	}
}
