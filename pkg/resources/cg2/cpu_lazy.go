package cg2

import (
	"io"
	"time"

	"github.com/ovsinc/resources-rate-limits/internal/utils"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"

	"go.uber.org/atomic"
)

var _ rescommon.Resourcer = (*CPUCG2Lazy)(nil)

type CPUCG2Lazy struct {
	ftotal io.ReadSeekCloser
	fused  io.ReadSeekCloser

	utilization *atomic.Float64
	tick        *time.Ticker
}

func NewCPULazy(
	done chan struct{},
	ftotal, fused io.ReadSeekCloser,
	dur time.Duration,
) (*CPUCG2Lazy, error) {
	if dur <= 0 {
		return nil, rescommon.ErrTickPeriodZero
	}

	cpu := &CPUCG2Lazy{
		ftotal:      ftotal,
		fused:       fused,
		utilization: &atomic.Float64{},
		tick:        time.NewTicker(dur),
	}

	cpu.init(done)

	// подождем для стабилизации 2 tick-периода + немного еще
	time.Sleep(2*dur + 100*time.Millisecond)

	return cpu, nil
}

func (cg *CPUCG2Lazy) Stop() {
	cg.ftotal.Close()
	cg.fused.Close()
}

func (cg *CPUCG2Lazy) Used() float64 {
	return cg.utilization.Load()
}

func (cg *CPUCG2Lazy) info() (total uint64, used uint64, err error) {
	return getCPUInfo(cg.ftotal, cg.fused)
}

func (cg *CPUCG2Lazy) init(done chan struct{}) {
	var (
		lastused  uint64
		lasttotal uint64
	)

	go func() {
		for {
			select {
			case <-done:
				return

			case <-cg.tick.C:
				total, used, err := cg.info()
				if err != nil {
					// неявный признак ошибки
					cg.utilization.Store(0)
					continue
				}

				// на первом круге (lasttotal == 0) пропускаем установку значения утилизации
				if lasttotal > 0 {
					cg.utilization.Store(utils.CPUPercent(lastused, used, lasttotal, total))
				}

				lastused = used
				lasttotal = total

			}
		}
	}()
}
