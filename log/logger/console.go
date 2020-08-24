package logger

import (
	"os"
)

type ConsoleOutputer struct {
}

func NewConsoleOutputer() (*ConsoleOutputer) {
	return &ConsoleOutputer{}
}

//直接输出控制台
func (c *ConsoleOutputer) Write(data *LogData) {
	_,_ = os.Stdout.Write([]byte(data.Bytes()))
}

func (c *ConsoleOutputer) Close() {
}
