package gorpheus

type Logger interface {
	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Error(format string, args ...interface{})
}

const (
	LoggerSQL = 0 << iota
	LoggerGorpheus
	LoggerDebug
)

const (
	LogLevelInfo = 0 << iota
	LogLevelDebug
	LogLevelError
)
