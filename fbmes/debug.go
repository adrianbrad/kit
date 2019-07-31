package fbmes

var logger DebugLogger

type DebugLogger interface {
	Debugf(string, ...interface{})
}

func debug(message string, args ...interface{}) {
	if logger != nil {
		logger.Debugf(message, args...)
	}
}

func SetDebugLogger(l DebugLogger) {
	logger = l
}
