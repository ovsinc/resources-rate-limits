package cg2

import (
	"io"
	"time"

	"github.com/ovsinc/errors"
	"github.com/ovsinc/resources-rate-limits/internal/utils"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"

	"go.uber.org/atomic"
)

type MemCG2Lazy struct {
	ftotal io.ReadSeekCloser
	fused  io.ReadSeekCloser
	used   *atomic.Float64
	tick   *time.Ticker
	dur    time.Duration
	done   chan struct{}
}

func NewMemLazy(
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

	ftotal := conf.File(rescommon.CGroup2MemLimitPath)
	if ftotal == nil {
		return nil, rescommon.ErrNoResourceReadFile.
			WithOptions(
				errors.AppendContextInfo("ftotal", rescommon.CGroup2MemLimitPath),
			)
	}

	fused := conf.File(rescommon.CGroup2MemUsagePath)
	if fused == nil {
		return nil, rescommon.ErrNoResourceReadFile.
			WithOptions(
				errors.AppendContextInfo("fused", rescommon.CGroup2MemUsagePath),
			)
	}

	mem := &MemCG2Lazy{
		dur:    dur,
		used:   &atomic.Float64{},
		tick:   time.NewTicker(dur),
		ftotal: ftotal,
		fused:  fused,
		done:   done,
	}

	mem.init()

	// подождем для стабилизации tick-период + немного еще
	time.Sleep(dur + (100 * time.Millisecond))

	return mem, nil
}

func (cg *MemCG2Lazy) Used() float64 {
	return cg.used.Load()
}

func (cg *MemCG2Lazy) Stop() {
	cg.tick.Stop()
}

func (cg *MemCG2Lazy) info() (uint64, uint64, error) {
	return getMemInfo(cg.ftotal, cg.fused)
}

func (cg *MemCG2Lazy) init() {
	cg.tick = time.NewTicker(cg.dur)
	go func() {
		for {
			select {
			case <-cg.done:
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
