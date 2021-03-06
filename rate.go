package resourcesratelimits

import (
	"time"

	"github.com/ovsinc/errors"
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
	Time           time.Time
}

// Limiter интерфейс лимитера.
type Limiter interface {
	// Основной метод
	Limit() *RateReply
	With(ops ...Option) Limiter
}

type resourceLimit struct {
	cpuRes, ramRes ResourceViewer
}

func New(ops ...Option) (Limiter, error) {
	rlp := new(resourceLimit)

	for _, op := range ops {
		op(rlp)
	}

	var err error

	// если не задано н одного ресорсера, устанавливаем автоматически
	if rlp.cpuRes == nil && rlp.ramRes == nil {
		rlp.cpuRes, rlp.ramRes, err = NewSimple()
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
	repl := &RateReply{
		Time: time.Now(),
	}

	if rl.ramRes != nil {
		repl.RAMUsed = rl.ramRes.Used()
	}

	if rl.cpuRes != nil {
		repl.CPUUtilization = rl.cpuRes.Used()
	}

	return repl
}

func (rl *resourceLimit) With(ops ...Option) Limiter {
	newrl := *rl
	for _, op := range ops {
		op(&newrl)
	}
	return &newrl
}
