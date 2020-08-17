package worker

import (
	"context"
	"ferry/pkg/logger"
	"os/exec"
	"syscall"

	"github.com/RichardKnop/machinery/v1/tasks"
)

var asyncTaskMap map[string]interface{}

func executeTaskBase(scriptPath string) {
	command := exec.Command("/bin/bash", "-c", scriptPath) //初始化Cmd
	err := command.Start()                                 //运行脚本
	if nil != err {
		logger.Errorf("task exec failed，%v", err.Error())
		return
	}

	logger.Info("Process PID:", command.Process.Pid)

	err = command.Wait() //等待执行完成
	if nil != err {
		logger.Errorf("task exec failed，%v", err.Error())
		return
	}

	logger.Info("ProcessState PID:", command.ProcessState.Pid())
	logger.Info("Exit Code", command.ProcessState.Sys().(syscall.WaitStatus).ExitStatus())
}

// ExecCommand 异步任务
func ExecCommand(classify string, scriptPath string) error {
	if classify == "shell" {
		logger.Info("start exec shell...", scriptPath)
		executeTaskBase(scriptPath)
		return nil
	} else if classify == "python" {
		logger.Info("start exec python...", scriptPath)
		executeTaskBase(scriptPath)
		return nil
	}
	return nil
}

func SendTask(ctx context.Context, classify string, scriptPath string) {
	args := make([]tasks.Arg, 0)
	args = append(args, tasks.Arg{
		Name:  "classify",
		Type:  "string",
		Value: classify,
	})
	args = append(args, tasks.Arg{
		Name:  "scriptPath",
		Type:  "string",
		Value: scriptPath,
	})
	task, _ := tasks.NewSignature("ExecCommandTask", args)
	task.RetryCount = 5
	_, err := AsyncTaskCenter.SendTaskWithContext(ctx, task)
	if err != nil {
		logger.Error(err.Error())
	}
}

func initAsyncTaskMap() {
	asyncTaskMap = make(map[string]interface{})
	asyncTaskMap["ExecCommandTask"] = ExecCommand
}
