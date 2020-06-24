package logOutputer

import (
	"context"
	"fmt"
	logBase "myRPC/log/base"
	"myRPC/trace"
	"path"
	"time"
)

var (
	lm                 *LoggerMgr
	DefaultLogChanSize = 10
	DefaultServiceName = "default"
	DefaultLogLevel    = "debug"
)

type LoggerMgr struct {
	outputers   []OutputerInterface
	chanSize    int
	level       string
	logDataChan chan *logBase.LogData
	serviceName string
}

func initLogger(level string, chanSize int, serviceName string) {
	lm = &LoggerMgr{
		chanSize:    chanSize,
		level:       level,
		serviceName: serviceName,
		logDataChan: make(chan *logBase.LogData, chanSize),
	}
	go lm.run()
}

func InitLogger(level string, chanSize int, serviceName string) {
	initLogger(level, chanSize, serviceName)
}

func SetLevel(level string) {
	lm.level = level
}

func (l *LoggerMgr) run() {
	for data := range l.logDataChan {
		for _, outputer := range l.outputers {
			outputer.Write(data)
		}
	}
}

func AddOutputer(ouputer OutputerInterface) {
	if lm == nil {
		initLogger(DefaultLogLevel, DefaultLogChanSize, DefaultServiceName)
	}
	lm.outputers = append(lm.outputers, ouputer)
	return
}

func Debug(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx, "debug", format, args...)
}

func Trace(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx, "trace", format, args...)
}

func Access(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx, "access", format, args...)
}

func Info(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx, "info", format, args...)
}

func Warn(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx, "warn", format, args...)
}

func Error(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx, "error", format, args...)
}

func Stop() {
	close(lm.logDataChan)
	for _, outputer := range lm.outputers {
		outputer.Close()
	}
	lm = nil
}

func writeLog(ctx context.Context, level string, format string, args ...interface{}) {
	if lm == nil {
		initLogger(DefaultLogLevel, DefaultLogChanSize, DefaultServiceName)
	}

	now := time.Now()
	nowStr := now.Format("2006-01-02 15:04:05.999")

	fileName, lineNo := logBase.GetLineInfo()
	fileName = path.Base(fileName)
	msg := fmt.Sprintf(format, args...)

	logData := &logBase.LogData{
		Message:     msg,
		CurTime:     now,
		TimeStr:     nowStr,
		Level:       level,
		Filename:    fileName,
		LineNo:      lineNo,
		TraceId:     trace.GetTraceId(ctx),
		ServiceName: lm.serviceName,
	}

	if level == "access" {
		fields := logBase.GetFields(ctx)
		if fields != nil {
			logData.Fields = fields
		}
	}

	select {
	case lm.logDataChan <- logData:
	default:
		return
	}
}
