package main

import (
	"ferry/cmd"
	"ferry/pkg/task"
)

func main() {
	// 启动异步任务队列
	go task.Start()

	cmd.Execute()
}
