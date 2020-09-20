package task

/*
  @Author : lanyulei
*/

import (
	"context"
	"ferry/pkg/task/worker"
)

func Send(classify string, scriptPath string, params string) {
	worker.SendTask(context.Background(), classify, scriptPath, params)
}
