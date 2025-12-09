package log

import "github.com/go-logr/logr"

func Debug(log logr.Logger, msg string, keysAndValues ...interface{}) {
	log.V(1).Info(msg, keysAndValues...)
}
