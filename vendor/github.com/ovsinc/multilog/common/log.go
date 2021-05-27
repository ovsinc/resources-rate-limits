package common

import "fmt"

// Logger интерфейс логгера.
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

func Format(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
