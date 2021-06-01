package cg1

import (
	"io"
	"time"

	"github.com/ovsinc/errors"
	"github.com/ovsinc/resources-rate-limits/internal/utils"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"

	"go.uber.org/atomic"
)

type MemCG1Lazy struct {
	ftotal io.ReadSeekCloser
	fused  io.ReadSeekCloser
	used   *atomic.Float64
	dur    time.Duration
	done   chan struct{}
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

	var err error

	ftotal := conf.File(rescommon.CGroupMemLimitPath)
	if ftotal == nil {
		err = errors.Wrap(
			err,
			rescommon.ErrNoResourceReadFile.WithOptions(
				errors.AppendContextInfo("ftotal", rescommon.CGroup2MemLimitPath),
				errors.AppendOperations("NewMemLazy"),
			),
		)
	}

	fused := conf.File(rescommon.CGroupMemUsagePath)
	if fused == nil {
		err = errors.Wrap(
			err,
			rescommon.ErrNoResourceReadFile.WithOptions(
				errors.AppendContextInfo("fused", rescommon.CGroup2MemUsagePath),
				errors.AppendOperations("NewMemLazy"),
			),
		)
	}

	if err != nil {
		return nil, err
	}

	mem := &MemCG1Lazy{
		dur:    dur,
		used:   &atomic.Float64{},
		ftotal: ftotal,
		fused:  fused,
		done:   done,
	}

	mem.init()

	return mem, nil
}

func (cg *MemCG1Lazy) Used() float64 {
	return cg.used.Load()
}

func (cg *MemCG1Lazy) info() (total, used uint64, err error) {
	_, err = cg.ftotal.Seek(0, 0)
	if err != nil {
		return 0, 0, err
	}

	_, err = cg.fused.Seek(0, 0)
	if err != nil {
		return 0, 0, err
	}

	total, err = readInfo(cg.ftotal)
	if err != nil {
		return 0, 0, err
	}

	used, err = readInfo(cg.fused)
	if err != nil {
		return 0, 0, err
	}

	return total, used, nil
}

func (cg *MemCG1Lazy) init() {
	tick := time.NewTicker(cg.dur)
	go func() {
		for {
			select {
			case <-cg.done:
				return
			case <-tick.C:
				total, used, err := cg.info()
				if err != nil {
					cg.used.Store(rescommon.FailValue)
					rescommon.DbgErrCommon("MemCG1Lazy", err)
					continue
				}
				p := utils.Percent(float64(used), float64(total))
				cg.used.Store(p)
				rescommon.DbgInfRAM("MemCG1Lazy", used, total, p)
			}
		}
	}()
}
