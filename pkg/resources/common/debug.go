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
	FmtErr     = "[%s]<ERR> check resource fails with %v"
	FmtCPUInfo = "[%s]<INFO> last: %d/%d now: %d/%d (%.2f%%)"
	FmtRAMInfo = "[%s]<INFO> now: %d/%d (%.2f%%)"
)

func DbgErrCommon(method string, err error) {
	Debug(FmtErr, method, err)
}

func DbgInfCPU(method string, lastused, used, lasttotal, total uint64, percents float64) {
	Debug(FmtCPUInfo, method, lastused, lasttotal, used, total, percents)
}

func DbgInfRAM(method string, used, total uint64, percents float64) {
	Debug(FmtRAMInfo, method, used, total, percents)
}
