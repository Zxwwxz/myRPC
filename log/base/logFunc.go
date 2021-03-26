package logBase
//不同等级的日志打印开放方法
func Debug(format string, args ...interface{}) {
	writeLog( "debug", format, args...)
}

func Info(format string, args ...interface{}) {
	writeLog( "info", format, args...)
}

func Warn(format string, args ...interface{}) {
	writeLog( "warn", format, args...)
}

func Fatal(format string, args ...interface{}) {
	writeLog( "fatal", format, args...)
}
