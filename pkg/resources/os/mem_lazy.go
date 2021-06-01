package os

import (
	"time"

	"github.com/ovsinc/errors"
	"github.com/ovsinc/resources-rate-limits/internal/utils"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"

	"go.uber.org/atomic"
)

type MemOSLazy struct {
	f    rescommon.ReadSeekCloser
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
				errors.AppendOperations("os.NewMemLazy"),
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

	return mem, nil
}

func (mem *MemOSLazy) Used() float64 {
	return mem.used.Load()
}

func (mem *MemOSLazy) info() (uint64, uint64, error) {
	return GetMemInfo(mem.f)
}

func (mem *MemOSLazy) init() {
	store := func() {
		total, used, err := mem.info()
		if err != nil {
			rescommon.DbgErrCommon("MemOSLazy", err)
			mem.used.Store(rescommon.FailValue)
			return
		}

		p := utils.Percent(float64(used), float64(total))
		mem.used.Store(p)
		rescommon.DbgInfRAM("MemOSLazy", used, total, p)
	}

	store()

	tick := time.NewTicker(mem.dur)

	go func() {
		for {
			select {
			case <-mem.done:
				mem.used.Store(rescommon.DoneValue)
				return

			case <-tick.C:
				store()
			}
		}
	}()
}
