package logger

type LoggerInterface interface {
	Write(data *LogData)error
	Close()error
}