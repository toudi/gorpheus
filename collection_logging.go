package gorpheus

func (c *Collection) SetLogger(id uint, logger Logger) {
	c.loggers[id] = logger
}

func (c *Collection) Log(loggerId uint, level uint, format string, args ...interface{}) {
	logger, exists := c.loggers[loggerId]
	if !exists {
		return
	}

	if level == LogLevelInfo {
		logger.Info(format, args...)
	} else if level == LogLevelDebug {
		logger.Debug(format, args...)
	} else if level == LogLevelError {
		logger.Error(format, args...)
	}
}
