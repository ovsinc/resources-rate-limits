package os

import (
	"io"
	"time"

	"github.com/ovsinc/errors"
	"github.com/ovsinc/resources-rate-limits/internal/utils"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"

	"go.uber.org/atomic"
)

type CPUOSLazy struct {
	dur         time.Duration
	f           io.ReadSeekCloser
	utilization *atomic.Float64
	tick        *time.Ticker
	done        chan struct{}
}

func NewCPULazy(
	done chan struct{},
	conf rescommon.ResourceConfiger,
	dur time.Duration,
) (rescommon.Resourcer, error) {
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

	// подождем для стабилизации 2 tick-периода + немного еще
	time.Sleep(2*dur + 100*time.Millisecond)

	return cpu, nil
}

func (cpu *CPUOSLazy) Used() float64 {
	return cpu.utilization.Load()
}

func (cpu *CPUOSLazy) Stop() {
	cpu.tick.Stop()
}

func (cpu *CPUOSLazy) info() (total uint64, used uint64, err error) {
	return getCPUInfo(cpu.f)
}

func (cpu *CPUOSLazy) init() {
	cpu.tick = time.NewTicker(cpu.dur)

	var (
		lastused  uint64
		lasttotal uint64
	)

	go func() {
		for {
			select {
			case <-cpu.done:
				return

			case <-cpu.tick.C:
				total, used, err := cpu.info()
				if err != nil {
					cpu.utilization.Store(0)
				}

				// на первом круге (lasttotal == 0) пропускаем установку значения утилизации
				if lasttotal > 0 {
					cpu.utilization.Store(utils.CPUPercent(lastused, used, lasttotal, total))
				}

				lastused = used
				lasttotal = total

			}
		}
	}()
}
