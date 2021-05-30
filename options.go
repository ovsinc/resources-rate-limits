package resourcesratelimits

import (
	"time"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
)

const (
	DefaultMemoryUsageBarrierPercentage    = 80.0
	DefaultCPUUtilizationBarrierPercentage = 80.0
	FailValue                              = rescommon.FailValue
)

type Option func(*resourceLimit)

var (
	AppendCPUResourcer = func(res ResourceViewer) Option {
		return func(rlp *resourceLimit) {
			rlp.cpuRes = res
		}
	}
	AppendRAMResourcer = func(res ResourceViewer) Option {
		return func(rlp *resourceLimit) {
			rlp.ramRes = res
		}
	}
	SetDebug = func(debug bool) Option {
		return func(rlp *resourceLimit) {
			rescommon.IsDebug = debug
		}
	}
)

func NewSimple() (resCPU, resRAM ResourceViewer, err error) {
	rt := Check()

	resCPU, err = AutoCPUSimple(rt.Type())
	if err != nil {
		return nil, nil, err
	}

	resRAM, err = AutoRAMSimple(rt.Type())
	if err != nil {
		return nil, nil, err
	}

	return resCPU, resRAM, nil
}

func MustNewSimple() (resCPU, resRAM ResourceViewer) {
	var err error
	resCPU, resRAM, err = NewSimple()
	if err != nil {
		panic(err)
	}
	return
}

func NewLazy() (resCPU, resRAM ResourceViewer, done chan struct{}, err error) {
	cpuDur := rescommon.DefaultDuration
	ramDur := rescommon.DefaultDuration + 111*time.Millisecond

	done = make(chan struct{})

	rt := Check()
	if err := rt.Init(); err != nil {
		return nil, nil, nil, err
	}

	go func() {
		<-done
		rt.Stop()
	}()

	resCPU, err = AutoLazyCPU(rt, done, cpuDur)
	if err != nil {
		close(done)
		return nil, nil, nil, err
	}

	resRAM, err = AutoLazyRAM(rt, done, ramDur)
	if err != nil {
		close(done)
		return nil, nil, nil, err
	}

	return resCPU, resRAM, done, nil
}

func MustNewLazy() (resCPU, resRAM ResourceViewer, done chan struct{}) {
	var err error
	resCPU, resRAM, done, err = NewLazy()
	if err != nil {
		panic(err)
	}
	return
}
