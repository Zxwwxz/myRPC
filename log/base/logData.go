package logBase

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

type LogData struct {
	//当前时间
	CurTime     time.Time
	//日志内容
	Message     string
	//时间格式
	TimeStr     string
	//日志等级
	Level       string
	//日志文件名
	Filename    string
	//日志行号
	LineNo      int
	//追踪id
	TraceId     string
	//服务名称
	ServiceName string
	//日志kv集合，统一输出
	Fields      *LogField
}

type KeyVal struct {
	key interface{}
	val interface{}
}

type LogField struct {
	kvs       []KeyVal
	fieldLock sync.Mutex
}

func writeField(buffer *bytes.Buffer, field, sep string) {
	buffer.WriteString(field)
	buffer.WriteString(sep)
}

//打印格式
func (l *LogData) Bytes() []byte {
	var buffer bytes.Buffer
	levelStr := l.Level
	writeField(&buffer, l.TimeStr, SpaceSep)
	writeField(&buffer, levelStr, SpaceSep)
	writeField(&buffer, l.ServiceName, SpaceSep)
	writeField(&buffer, l.Filename, ColonSep)
	writeField(&buffer, fmt.Sprintf("%d", l.LineNo), SpaceSep)
	writeField(&buffer, l.TraceId, SpaceSep)
	//access统一输出
	if l.Level == "access" && l.Fields != nil {
		for _, field := range l.Fields.kvs {
			writeField(&buffer, fmt.Sprintf("%v=%v", field.key, field.val), SpaceSep)
		}
	}
	writeField(&buffer, l.Message, LineSep)
	return buffer.Bytes()
}

func GetLineInfo() (fileName string, lineNo int) {
	_, fileName, lineNo, _ = runtime.Caller(3)
	return
}

func (l *LogField) AddField(key, val interface{}) {
	l.fieldLock.Lock()
	l.kvs = append(l.kvs, KeyVal{key: key, val: val})
	l.fieldLock.Unlock()
}

type kvsIdKey struct{}

func WithFieldContext(ctx context.Context) context.Context {
	fields := &LogField{}
	return context.WithValue(ctx, kvsIdKey{}, fields)
}

//添加日志到field
func AddField(ctx context.Context, key string, val interface{}) {
	field := GetFields(ctx)
	if field == nil {
		return
	}
	field.AddField(key, val)
}

func GetFields(ctx context.Context) *LogField {
	field, ok := ctx.Value(kvsIdKey{}).(*LogField)
	if !ok {
		return nil
	}
	return field
}
