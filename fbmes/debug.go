package fbmes

var DEBUG bool
var Logger DebugLogger

type DebugLogger interface {
	Debugf(string, ...interface{})
}

func Debug(message string, args ...interface{}) {
	if DEBUG {
		if Logger != nil {
			Logger.Debugf(message, args...)
		}
	}
}
