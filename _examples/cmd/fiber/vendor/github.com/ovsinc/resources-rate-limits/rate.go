package resourcesratelimits

import (
	"time"

	"github.com/ovsinc/errors"

	"github.com/ovsinc/resources-rate-limits/pkg/resources"
)

var (
	ErrCPUUtilizationIsZero = errors.New("CPU utilization is zero")
	ErrRAMUtilizationIsZero = errors.New("RAM utilization is zero")

	ErrNoResourcer = errors.New("not any resourcers set")
)

var _ Limiter = (*resourceLimit)(nil)

type RateReply struct {
	RAMUsed        float64
	CPUUtilization float64
	Time           *time.Time
	Err            error
}

type Limiter interface {
	Limit() *RateReply
	Shutdown()
}

type resourceLimit struct {
	cpuRes, ramRes resources.Resourcer
	conf           *RateLimitConfig
}

func New(ops ...Option) (Limiter, error) {
	rlp := new(resourceLimit)

	for _, op := range ops {
		op(rlp)
	}

	// конфиг всегда должен быть
	if rlp.conf == nil {
		rlp.conf = DefaultRateLimitConfig
	}

	var err error

	// если не задано н одного ресорсера, устанавливаем автоматически
	if rlp.cpuRes == nil && rlp.ramRes == nil {
		rlp.cpuRes, err = resources.AutoCPU()
		if err != nil {
			return nil, err
		}

		rlp.ramRes, err = resources.AutoRAM()
		if err != nil {
			return nil, err
		}
	}

	return rlp, nil
}

func MustNew(ops ...Option) Limiter {
	l, err := New(ops...)
	if err != nil {
		print(err)
	}
	return l
}

func (rl *resourceLimit) Limit() *RateReply {
	repl := new(RateReply)

	if rl.ramRes != nil {
		repl.RAMUsed = rl.ramRes.Used()
	}

	if rl.cpuRes != nil {
		repl.CPUUtilization = rl.cpuRes.Used()
	}

	switch {
	case repl.RAMUsed == 0.0:
		repl.Err = ErrRAMUtilizationIsZero
		if rl.conf.ErrorHandler != nil {
			rl.conf.ErrorHandler(rl.conf, repl.Err)
		}

	case repl.CPUUtilization == 0.0:
		if rl.conf.ErrorHandler != nil {
			rl.conf.ErrorHandler(rl.conf, repl.Err)
		}

		repl.Err = ErrCPUUtilizationIsZero
		now := time.Now()
		repl.Time = &now

	case repl.RAMUsed >= rl.conf.MemoryUsageBarrierPercentage,
		repl.CPUUtilization >= rl.conf.CPUUtilizationBarrierPercentage:
		if rl.conf.LimitHandler != nil {
			rl.conf.LimitHandler(rl.conf)
		}

		now := time.Now()
		repl.Time = &now
	}

	return repl
}

func (rl *resourceLimit) Shutdown() {
	if rl.cpuRes != nil {
		rl.cpuRes.Stop()
	}
	if rl.ramRes != nil {
		rl.ramRes.Stop()
	}
}
