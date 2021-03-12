package service

import (
	"ferry/pkg/task"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

/*
  @Author : lanyulei
*/

func ExecTask(taskList []string, params string) {
	for _, taskName := range taskList {
		filePath := fmt.Sprintf("%v/%v", viper.GetString("script.path"), taskName)
		if strings.HasSuffix(filePath, ".py") {
			task.Send("python", filePath, params)
		} else if strings.HasSuffix(filePath, ".sh") {
			task.Send("shell", filePath, params)
		}
	}
}
