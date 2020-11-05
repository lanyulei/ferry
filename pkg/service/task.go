package service

import (
	"ferry/pkg/task"
	"ferry/tools/config"
	"path/filepath"
	"strings"
)

/*
  @Author : lanyulei
*/

func ExecTask(taskList []string, params string) {
	for _, taskName := range taskList {
		filePath := filepath.Join(config.ScriptPath, taskName)
		if strings.HasSuffix(filePath, ".py") {
			task.Send("python", filePath, params)
		} else if strings.HasSuffix(filePath, ".sh") {
			task.Send("shell", filePath, params)
		}
	}
}
