package common

import "log"

var IsDebug = false

func Debug(format string, arg ...interface{}) {
	if IsDebug {
		log.Panicf(format, arg...)
	}
}
