package os

import (
	"io"
	"time"

	"github.com/ovsinc/resources-rate-limits/internal/utils"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"

	"go.uber.org/atomic"
)

var _ rescommon.Resourcer = (*CPUOSLazy)(nil)

type CPUOSLazy struct {
	f           io.ReadSeekCloser
	utilization *atomic.Float64
	tick        *time.Ticker
}

func NewCPULazy(done chan struct{}, f io.ReadSeekCloser, dur time.Duration) (*CPUOSLazy, error) {
	if dur <= 0 {
		return nil, rescommon.ErrTickPeriodZero
	}

	cpu := &CPUOSLazy{
		f:           f,
		utilization: &atomic.Float64{},
		tick:        time.NewTicker(dur),
	}

	err := cpu.init(done)
	if err != nil {
		return nil, err
	}

	// подождем для стабилизации 2 tick-периода + немного еще
	time.Sleep(2*dur + 100*time.Millisecond)

	return cpu, nil
}

func (cpu *CPUOSLazy) Used() float64 {
	return cpu.utilization.Load()
}

func (cpu *CPUOSLazy) Stop() {
	cpu.f.Close()
	cpu.tick.Stop()
}

func (cpu *CPUOSLazy) info() (total uint64, used uint64, err error) {
	return getCPUInfo(cpu.f)
}

func (cpu *CPUOSLazy) init(done chan struct{}) error {
	var (
		errGlob   atomic.Error
		lastused  uint64
		lasttotal uint64
	)

	go func() {
		for {
			select {
			case <-done:
				return

			case <-cpu.tick.C:
				total, used, err := cpu.info()
				if err != nil {
					errGlob.Store(err)
					return
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

	return errGlob.Load()
}
