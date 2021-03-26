package logBase

import (
	"errors"
	"fmt"
	"myRPC/log/logger"
	"path"
	"time"
)

//日志等级
const (
	log_type_debug = "debug"
	log_type_info = "info"
	log_type_warn = "warn"
	log_type_fatal = "fatal"
)

const (
	default_chan_size = 100
)

var (
	logMgr                 *LogMgr
)

//日志管理器
type LogMgr struct {
	//开关
	switchOn    bool
	//每个等级字符转数字
	logType     map[string]int
	//每个等级日志对象
	loggers     map[string]logger.LoggerInterface
	//日志通道大小
	chanSize    int
	//当前等级
	level       string
	//日志打印通道
	logDataChan chan *logger.LogData
}

//服务创建时初始化日志管理器
func InitLogger(switchOn bool,level string, chanSize int,params map[interface{}]interface{}) {
	logMgr = &LogMgr{
		switchOn:    false,
		chanSize:    default_chan_size,
		level:       log_type_debug,
		logDataChan: make(chan *logger.LogData, chanSize),
	}
	if chanSize != 0 {
		logMgr.chanSize = chanSize
	}
	logMgr.switchOn = switchOn
	//设置日志等级
	_ = logMgr.SetLevel(level)
	logMgr.logType = map[string]int{log_type_debug:1,log_type_info:2,log_type_warn:3,log_type_fatal:4}
	logMgr.loggers = make(map[string]logger.LoggerInterface,len(logMgr.logType))
	//新增日志等级对象
	logMgr.addLogger(params)
	go logMgr.run()
}

//服务停止时关闭日志管理器
func Stop() {
	close(logMgr.logDataChan)
	for _, tempLogger := range logMgr.loggers {
		_ = tempLogger.Close()
	}
	logMgr = nil
}

func GetLogMgr()(*LogMgr) {
	return logMgr
}

func (l *LogMgr)SetSwitchOn(switchOn bool)(error) {
	if l == nil {
		return errors.New("LogMgr nil")
	}
	l.switchOn = switchOn
	return nil
}

//设置日志等级
func (l *LogMgr)SetLevel(level string)(error) {
	if l == nil {
		return errors.New("LogMgr nil")
	}
	if _,ok := l.logType[level];ok == false{
		return errors.New("level illegal")
	}
	l.level = level
	return nil
}

//添加不同等级的日志打印对象
func (l *LogMgr)addLogger(params map[interface{}]interface{}) {
	//如果想做不同等级的日志合并，那就这里不要创建多个对象，所有等级指向通过对象
	for logType,_ := range logMgr.logType{
		tempLogger,err := logger.NewFileOutputer(params,logType)
		if err != nil{
			continue
		}
		logMgr.loggers[logType] = tempLogger
	}
	return
}

//轮询通道写日志
func (l *LogMgr) run() {
	for data := range l.logDataChan {
		if tempLogger,ok := l.loggers[data.GetLevel()];ok == true{
			_ = tempLogger.Write(data)
		}
	}
}

func (l *LogMgr) Stop() {

}

//将日志写入通道中
func writeLog(level string, format string, args ...interface{}) {
	if logMgr == nil || logMgr.switchOn == false{
		return
	}
	//日志等级不足以打印
	reqLevel := logMgr.logType[level]
	mgrLevel := logMgr.logType[logMgr.level]
	if reqLevel < mgrLevel {
		return
	}
	//获取正在执行的文件名和行号
	fileName, lineNo := logger.GetLineInfo()
	fileName = path.Base(fileName)
	msg := fmt.Sprintf(format, args...)
	//创建日志内容对象，放到通道中
	logData := logger.NewLogData(msg,level,fileName,lineNo,time.Now())
	select {
	case logMgr.logDataChan <- logData:
	default:
		return
	}
}
