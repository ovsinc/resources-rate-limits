package common

import (
	"log"
	"os"

	"github.com/ovsinc/multilog/golog"
)

var (
	IsDebug  = false
	Debugger = golog.New(log.New(os.Stdout, "ovsinc/resources-rate-limits", log.LstdFlags))
)

func Debug(format string, arg ...interface{}) {
	if IsDebug {
		Debugger.Debugf(format, arg...)
	}
}

const (
	FmtErr          = "[%s]<ERR> check resource fails with %v"
	FmtCPUInfo      = "[%s]<INFO> last: %d/%d now: %d/%d"
	FmtCPUFirstInfo = "[%s]<INFO> first loop last: %d/%d"
	FmtRAMInfo      = "[%s]<INFO> now: %d/%d"
)

func DbgErrCommon(method string, err error) {
	Debug(FmtErr, method, err)
}

func DbgInfCPU(method string, lastused, used, lasttotal, total uint64) {
	Debug(FmtCPUInfo, method, lastused, lasttotal, used, total)
}

func DbgInfRAM(method string, used, total uint64) {
	Debug(FmtRAMInfo, method, used, total)
}

func DbgInfCPUFirst(method string, lastused, lasttotal uint64) {
	Debug(FmtCPUFirstInfo, method, lastused, lasttotal)
}
