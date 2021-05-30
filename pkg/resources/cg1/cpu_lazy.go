package cg1

import (
	"io"
	"time"

	"github.com/ovsinc/errors"
	"github.com/ovsinc/resources-rate-limits/internal/utils"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"

	"go.uber.org/atomic"
)

type CPUCG1Lazy struct {
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
				errors.AppendContextInfo("ftotal", rescommon.CGroup2CPULimitPath),
				errors.AppendOperations("cg1.NewCPULazy"),
			)
	}

	fused := conf.File(rescommon.CGroup2CPUUsagePath)
	if fused == nil {
		return nil, rescommon.ErrNoResourceReadFile.
			WithOptions(
				errors.AppendOperations("cg1.NewCPULazy"),
				errors.AppendContextInfo("fused", rescommon.CGroup2CPUUsagePath),
			)
	}

	cpu := &CPUCG1Lazy{
		dur:         dur,
		done:        done,
		ftotal:      ftotal,
		fused:       fused,
		utilization: &atomic.Float64{},
	}

	cpu.init()

	return cpu, nil
}

func (cg *CPUCG1Lazy) Used() float64 {
	return cg.utilization.Load()
}

func (cg *CPUCG1Lazy) info() (total uint64, used uint64, err error) {
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

func (cg *CPUCG1Lazy) init() {
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
					rescommon.DbgErrCommon("CPUCG1Lazy", err)
					continue
				}

				// на первом круге (lasttotal == 0) пропускаем установку значения утилизации
				if lasttotal > 0 {
					cg.utilization.Store(utils.CPUPercent(lastused, used, lasttotal, total))
					rescommon.DbgInfCPU("CPUCG1Lazy", lastused, used, lasttotal, total)
				}

				lastused = used
				lasttotal = total

				rescommon.DbgInfCPUFirst("CPUCG1Lazy", lastused, lasttotal)
			}
		}
	}()
}