package cg2

import (
	"io"
	"time"

	"github.com/ovsinc/errors"
	"github.com/ovsinc/resources-rate-limits/internal/utils"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"

	"go.uber.org/atomic"
)

type CPUCG2Lazy struct {
	dur         time.Duration
	ftotal      io.ReadSeekCloser
	fused       io.ReadSeekCloser
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

	ftotal := conf.File(rescommon.CGroup2CPULimitPath)
	if ftotal == nil {
		return nil, rescommon.ErrNoResourceReadFile.
			WithOptions(
				errors.AppendOperations("cg2.NewCPULazy"),
				errors.AppendContextInfo("ftotal", rescommon.CGroup2CPULimitPath),
			)
	}

	fused := conf.File(rescommon.CGroup2CPUUsagePath)
	if fused == nil {
		return nil, rescommon.ErrNoResourceReadFile.
			WithOptions(
				errors.AppendOperations("cg2.NewCPULazy"),
				errors.AppendContextInfo("fused", rescommon.CGroup2CPUUsagePath),
			)
	}

	cpu := &CPUCG2Lazy{
		dur:         dur,
		done:        done,
		ftotal:      ftotal,
		fused:       fused,
		utilization: &atomic.Float64{},
	}

	cpu.init()

	return cpu, nil
}

func (cg *CPUCG2Lazy) Used() float64 {
	return cg.utilization.Load()
}

func (cg *CPUCG2Lazy) info() (total uint64, used uint64, err error) {
	return getCPUInfo(cg.ftotal, cg.fused)
}

func (cg *CPUCG2Lazy) init() {
	var (
		lastused  uint64
		lasttotal uint64
	)

	tick := time.NewTicker(cg.dur)

	go func() {
		for {
			select {
			case <-cg.done:
				return

			case <-tick.C:
				total, used, err := cg.info()
				if err != nil {
					// неявный признак ошибки
					cg.utilization.Store(rescommon.FailValue)
					rescommon.DbgErrCommon("CPUCG2Lazy", err)
					continue
				}

				// на первом круге (lasttotal == 0) пропускаем установку значения утилизации
				if lasttotal > 0 {
					p := utils.CPUPercent(lastused, used, lasttotal, total)
					cg.utilization.Store(p)
					rescommon.DbgInfCPU("CPUCG2Lazy", lastused, used, lasttotal, total, p)
				}

				lastused = used
				lasttotal = total
			}
		}
	}()
}
