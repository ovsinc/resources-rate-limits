package os

import (
	"io"
	"time"

	"github.com/ovsinc/errors"
	"github.com/ovsinc/resources-rate-limits/internal/utils"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"

	"go.uber.org/atomic"
)

type MemOSLazy struct {
	f    io.ReadSeekCloser
	used *atomic.Float64
	done chan struct{}
	dur  time.Duration
}

func NewMemLazy(
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

	f := conf.File(rescommon.RAMFilenameInfoProc)
	if f == nil {
		return nil, rescommon.ErrNoResourceReadFile.
			WithOptions(
				errors.AppendContextInfo("f", rescommon.RAMFilenameInfoProc),
			)
	}

	mem := &MemOSLazy{
		f:    f,
		dur:  dur,
		used: &atomic.Float64{},
		done: done,
	}

	mem.init()

	// подождем для стабилизации tick-период + немного еще
	time.Sleep(dur + (100 * time.Millisecond))

	return mem, nil
}

func (mem *MemOSLazy) Used() float64 {
	return mem.used.Load()
}

func (mem *MemOSLazy) info() (uint64, uint64, error) {
	return getMemInfo(mem.f)
}

func (mem *MemOSLazy) init() {
	tick := time.NewTicker(mem.dur)

	go func() {
		for {
			select {
			case <-mem.done:
				return
			case <-tick.C:
				total, used, err := mem.info()
				if err != nil {
					rescommon.Debug("[MemOSLazy]<ERR> Check resource fails with %v", err)
					mem.used.Store(rescommon.FailValue)
				}

				mem.used.Store(utils.Percent(float64(used), float64(total)))
				rescommon.Debug(
					"[MemOSLazy]<INFO> now: %d/%d",
					used, total,
				)
			}
		}
	}()
}
