package logOutputer

import (
	logBase "myRPC/log/base"
	"os"
)

type ConsoleOutputer struct {
}

func NewConsoleOutputer() (log OutputerInterface) {
	log = &ConsoleOutputer{}
	return
}

//直接输出控制台
func (c *ConsoleOutputer) Write(data *logBase.LogData) {
	_,_ = os.Stdout.Write([]byte(data.Bytes()))
}

func (c *ConsoleOutputer) Close() {
}
