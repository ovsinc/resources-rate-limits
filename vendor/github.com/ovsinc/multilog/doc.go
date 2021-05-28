// Package multilog is a simple logging wrapper for common logging applications.
// The following loggers are supported:
// * logrus (https://pkg.go.dev/github.com/sirupsen/logrus)
// * golog (https://pkg.go.dev/log)
// * log15 (https://pkg.go.dev/github.com/inconshreveable/log15)
// * journald (https://pkg.go.dev/github.com/coreos/go-systemd/journal)
// * syslog (https://pkg.go.dev/log/syslog)
// * zap (https://pkg.go.dev/go.uber.org/zap)
//
// It is possible to combine supported loggers in a chain.
//
//
// RU:
// Простая обёртка для распространенных систем логгирования.
// Поддерживаются следующие логгеры:
// * logrus (https://pkg.go.dev/github.com/sirupsen/logrus)
// * golog (https://pkg.go.dev/log)
// * log15 (https://pkg.go.dev/github.com/inconshreveable/log15)
// * journald (https://pkg.go.dev/github.com/coreos/go-systemd/journal)
// * syslog (https://pkg.go.dev/log/syslog)
// * zap (https://pkg.go.dev/go.uber.org/zap)
//
// Возможно объединение логгеров в цепочку.
//
//
package multilog // import "github.com/ovsinc/multilog"
