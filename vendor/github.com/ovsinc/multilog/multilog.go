package multilog

import (
	pkglog "log"
	"os"

	"github.com/ovsinc/multilog/common"
	"github.com/ovsinc/multilog/golog"
)

// Logger logger interface.
// Logger интерфейс логгера.
type Logger common.Logger

// DefaultLogger sets default logger is used.
//
// DefaultLogger логгер, используемый по умолчанию.
// Можно переопределить.
var DefaultLogger Logger = golog.New(pkglog.New(os.Stderr, "", pkglog.LstdFlags))

// Debugf logs an message with Debug level.
// The default logger is used.
//
// Debugf залоггирует сообщение уровня Debug.
// Используется дефолтный логгер.
func Debugf(format string, args ...interface{}) {
	DefaultLogger.Debugf(format, args...)
}

// Infof logs an message with Info level.
// The default logger is used.
//
// Infof залоггирует сообщение уровня Info.
// Используется дефолтный логгер.
func Infof(format string, args ...interface{}) {
	DefaultLogger.Infof(format, args...)
}

// Warnf logs an message with Warning level.
// The default logger is used.
//
// Warnf залоггирует сообщение уровня Warning.
// Используется дефолтный логгер.
func Warnf(format string, args ...interface{}) {
	DefaultLogger.Warnf(format, args...)
}

// Errorf logs an message with Error level.
// The default logger is used.
//
// Errorf залоггирует сообщение уровня Error.
// Используется дефолтный логгер.
func Errorf(format string, args ...interface{}) {
	DefaultLogger.Errorf(format, args...)
}
