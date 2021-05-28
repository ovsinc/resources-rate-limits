package cg2

import (
	"io"
	"time"

	"github.com/ovsinc/resources-rate-limits/internal/utils"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"

	"go.uber.org/atomic"
)

var _ rescommon.Resourcer = (*MemCG2Lazy)(nil)

type MemCG2Lazy struct {
	ftotal io.ReadSeekCloser
	fused  io.ReadSeekCloser

	used *atomic.Float64
	tick *time.Ticker
}

func NewMemLazy(
	done chan struct{},
	ftotal, fused io.ReadSeekCloser,
	dur time.Duration,
) (*MemCG2Lazy, error) {
	if dur <= 0 {
		return nil, rescommon.ErrTickPeriodZero
	}

	mem := &MemCG2Lazy{
		used:   &atomic.Float64{},
		tick:   time.NewTicker(dur),
		ftotal: ftotal,
		fused:  fused,
	}

	mem.init(done)

	// подождем для стабилизации tick-период + немного еще
	time.Sleep(dur + (100 * time.Millisecond))

	return mem, nil
}

func (cg *MemCG2Lazy) Used() float64 {
	return cg.used.Load()
}

func (cg *MemCG2Lazy) Stop() {
	cg.ftotal.Close()
	cg.fused.Close()
}

func (cg *MemCG2Lazy) info() (uint64, uint64, error) {
	return getMemInfo(cg.ftotal, cg.fused)
}

func (cg *MemCG2Lazy) init(done chan struct{}) {
	go func() {
		for {
			select {
			case <-done:
				return
			case <-cg.tick.C:
				total, used, err := cg.info()
				if err != nil {
					cg.used.Store(0)
					continue
				}

				cg.used.Store(utils.Percent(float64(used), float64(total)))
			}
		}
	}()
}
