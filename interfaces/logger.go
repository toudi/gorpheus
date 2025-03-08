package interfaces

import "log/slog"

type LogMessage struct {
	Level   slog.Level
	Message string
	Payload []any
}

type Logger interface {
	Log(m LogMessage)
}

// DevNullLogger will be the default logger inside gorpheus
// just to avoid checking if logger != nil since that's error-prone.
type DevNullLogger struct{}

func (dnl *DevNullLogger) Log(m LogMessage) {}
