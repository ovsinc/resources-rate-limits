// Package golog implements the standard golang logger.
//
//
// RU:
// Package golog реализует стандартный логгер golang.
package golog

import (
	gosystemlog "log"

	log "github.com/ovsinc/multilog/common"
)

// New constructor of a logger that wraps the original logger.
//
//
// RU:
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
	l.logger.Printf(format, args...)
}

func (l *systemlog) Infof(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *systemlog) Warnf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *systemlog) Errorf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}
