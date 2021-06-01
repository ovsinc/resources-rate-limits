package os

import (
	"time"

	"github.com/ovsinc/errors"
	"github.com/ovsinc/resources-rate-limits/internal/utils"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"

	"go.uber.org/atomic"
)

type CPUOSLazy struct {
	dur         time.Duration
	f           rescommon.ReadSeekCloser
	utilization *atomic.Float64
	done        chan struct{}
}

func NewCPULazy(
	done chan struct{},
	conf rescommon.ResourceConfiger,
	dur time.Duration,
) (rescommon.ResourceViewer, error) {
	if dur <= 0 {
		return nil, rescommon.ErrTickPeriodZero
	}

	if conf == nil {
		return nil, rescommon.ErrNoResourceConfig
	}

	f := conf.File(rescommon.CPUfilenameInfoProc)
	if f == nil {
		return nil, rescommon.ErrNoResourceReadFile.
			WithOptions(
				errors.AppendOperations("os.NewCPULazy"),
				errors.AppendContextInfo("f", rescommon.CPUfilenameInfoProc),
			)
	}

	cpu := &CPUOSLazy{
		f:           f,
		utilization: &atomic.Float64{},
		dur:         dur,
		done:        done,
	}

	cpu.init()

	return cpu, nil
}

func (cpu *CPUOSLazy) Used() float64 {
	return cpu.utilization.Load()
}

func (cpu *CPUOSLazy) info() (total uint64, used uint64, err error) {
	return getCPUInfo(cpu.f)
}

func (cpu *CPUOSLazy) init() {
	tick := time.NewTicker(cpu.dur)

	var (
		lastused  uint64
		lasttotal uint64
	)

	go func() {
		for {
			select {
			case <-cpu.done:
				cpu.utilization.Store(rescommon.DoneValue)
				return

			case <-tick.C:
				total, used, err := cpu.info()
				if err != nil {
					rescommon.DbgErrCommon("CPUOSLazy", err)
					cpu.utilization.Store(rescommon.FailValue)
				}

				// на первом круге (lasttotal == 0) пропускаем установку значения утилизации
				if lasttotal > 0 {
					p := utils.CPUPercent(lastused, used, lasttotal, total)
					cpu.utilization.Store(p)
					rescommon.DbgInfCPU("CPUOSLazy", lastused, used, lasttotal, total, p)
				}

				lastused = used
				lasttotal = total
			}
		}
	}()
}
