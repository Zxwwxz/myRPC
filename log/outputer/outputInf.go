package logOutputer

import logBase "myRPC/log/base"

type OutputerInterface interface {
	Write(data *logBase.LogData)
	Close()
}