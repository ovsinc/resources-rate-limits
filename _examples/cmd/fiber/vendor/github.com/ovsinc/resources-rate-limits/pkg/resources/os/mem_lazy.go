package os

import (
	"io"
	"time"

	"github.com/ovsinc/resources-rate-limits/internal/utils"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"

	"go.uber.org/atomic"
)

var _ rescommon.Resourcer = (*MemOSLazy)(nil)

type MemOSLazy struct {
	f    io.ReadSeekCloser
	used *atomic.Float64
	tick *time.Ticker
}

func NewMemLazy(done chan struct{}, f io.ReadSeekCloser, dur time.Duration) (*MemOSLazy, error) {
	if dur <= 0 {
		return nil, rescommon.ErrTickPeriodZero
	}

	mem := &MemOSLazy{
		f:    f,
		used: &atomic.Float64{},
		tick: time.NewTicker(dur),
	}

	mem.init(done)

	// подождем для стабилизации tick-период + немного еще
	time.Sleep(dur + (100 * time.Millisecond))

	return mem, nil
}

func (mem *MemOSLazy) Used() float64 {
	return mem.used.Load()
}

func (mem *MemOSLazy) Stop() {
	mem.f.Close()
	mem.tick.Stop()
}

func (mem *MemOSLazy) info() (uint64, uint64, error) {
	return getMemInfo(mem.f)
}

func (mem *MemOSLazy) init(done chan struct{}) error {
	var errGlob atomic.Error

	go func() {
		for {
			select {
			case <-done:
				return
			case <-mem.tick.C:
				total, used, err := mem.info()
				if err != nil {
					errGlob.Store(err)
					return
				}

				mem.used.Store(utils.Percent(float64(used), float64(total)))
			}
		}
	}()

	return errGlob.Load()
}
