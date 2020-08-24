package logger

import (
	"bytes"
	"fmt"
	"runtime"
	"time"
)

const (
	SpaceSep           = " "
	ColonSep           = ":"
	LineSep            = "\n"
)

type LogData struct {
	//日志等级
	level       string
	//当前时间
	curTime     time.Time
	//文件名
	filename    string
	//行号
	lineNo      int
	//日志内容
	message     string
}

func NewLogData(message,level,filename string,lineNo int,curTime time.Time)(*LogData) {
	return &LogData{
		level:level,
		curTime:curTime,
		filename:filename,
		lineNo:lineNo,
		message:message,
	}
}

func (l *LogData)GetLevel()string {
	return l.level
}

func (l *LogData)writeField(buffer *bytes.Buffer, field, sep string) {
	buffer.WriteString(field)
	buffer.WriteString(sep)
}

//[2020-08-20 08:00:00:000][堆栈]内容
func (l *LogData) Bytes() []byte {
	var buffer bytes.Buffer
	l.writeField(&buffer, fmt.Sprintf("[%s]",l.level), SpaceSep)
	l.writeField(&buffer, fmt.Sprintf("[%s:%d]",l.curTime.Format("2006-01-01 15:04:05"),l.curTime.Nanosecond()/1000000), SpaceSep)
	l.writeField(&buffer, fmt.Sprintf("[%s:%d]",l.filename,l.lineNo), SpaceSep)
	l.writeField(&buffer, l.message, LineSep)
	return buffer.Bytes()
}

func GetLineInfo() (fileName string, lineNo int) {
	_, fileName, lineNo, _ = runtime.Caller(3)
	return
}