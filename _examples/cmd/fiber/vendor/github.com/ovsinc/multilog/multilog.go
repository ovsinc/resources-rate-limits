package multilog

import (
	pkglog "log"
	"os"

	"github.com/ovsinc/multilog/common"
	"github.com/ovsinc/multilog/golog"
)

type Logger common.Logger

// DefaultLogger логгер, используемый по умолчанию.
// Можно переопределить.
var DefaultLogger common.Logger = golog.New(pkglog.New(os.Stderr, "", pkglog.LstdFlags))

func Debugf(format string, args ...interface{}) {
	DefaultLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	DefaultLogger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	DefaultLogger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	DefaultLogger.Errorf(format, args...)
}
