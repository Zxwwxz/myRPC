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
	CurTime     time.Time
	Message     string
	TimeStr     string
	Level       string
	Filename    string
	LineNo      int
	TraceId     string
	ServiceName string
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

func (l *LogData) Bytes() []byte {
	var buffer bytes.Buffer
	levelStr := l.Level
	writeField(&buffer, l.TimeStr, SpaceSep)
	writeField(&buffer, levelStr, SpaceSep)
	writeField(&buffer, l.ServiceName, SpaceSep)
	writeField(&buffer, l.Filename, ColonSep)
	writeField(&buffer, fmt.Sprintf("%d", l.LineNo), SpaceSep)
	writeField(&buffer, l.TraceId, SpaceSep)
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
