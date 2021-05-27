// Package golog реализует стандартный логгер golang.
package golog

import (
	gosystemlog "log"

	log "github.com/ovsinc/multilog/common"
)

// New конструтор интерфейс для использования логгера log из состава golang.
// Оборачивает стандартный логгер l.
func New(l *gosystemlog.Logger) log.Logger {
	return &systemlog{
		logger: l,
	}
}

type systemlog struct {
	logger *gosystemlog.Logger
}

func (l *systemlog) Debugf(format string, args ...interface{}) {
	l.logger.Printf("DEBUG: "+format, args...)
}

func (l *systemlog) Infof(format string, args ...interface{}) {
	l.logger.Printf("INFO: "+format, args...)
}

func (l *systemlog) Warnf(format string, args ...interface{}) {
	l.logger.Printf("WARN: "+format, args...)
}

func (l *systemlog) Errorf(format string, args ...interface{}) {
	l.logger.Printf("ERR: "+format, args...)
}
