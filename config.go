package resourcesratelimits

import (
	"time"

	"github.com/ovsinc/resources-rate-limits/pkg/resources"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
)

const (
	DefaultMemoryUsageBarrierPercentage    = 80.0
	DefaultCPUUtilizationBarrierPercentage = 80.0
	FailValue                              = rescommon.FailValue
)

type Option func(*resourceLimit)

var (
	AppendCPUResourcer = func(res resources.ResourceViewer) Option {
		return func(rlp *resourceLimit) {
			rlp.cpuRes = res
		}
	}
	AppendRAMResourcer = func(res resources.ResourceViewer) Option {
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

func NewSimple() (resCPU, resRAM resources.ResourceViewer, err error) {
	resCPU, err = resources.AutoCPUSimple()
	if err != nil {
		return nil, nil, err
	}

	resRAM, err = resources.AutoRAMSimple()
	if err != nil {
		return nil, nil, err
	}

	return resCPU, resRAM, nil
}

func MustNewSimple() (resCPU, resRAM resources.ResourceViewer) {
	var err error
	resCPU, resRAM, err = NewSimple()
	if err != nil {
		panic(err)
	}
	return
}

func NewLazy() (resCPU, resRAM resources.ResourceViewer, done chan struct{}, err error) {
	cpuDur := rescommon.DefaultDuration
	ramDur := rescommon.DefaultDuration + 111*time.Millisecond

	done = make(chan struct{})

	resCPU, err = resources.AutoLazyCPU(done, cpuDur)
	if err != nil {
		return nil, nil, nil, err
	}

	resRAM, err = resources.AutoLazyRAM(done, ramDur)
	if err != nil {
		return nil, nil, nil, err
	}

	return resCPU, resRAM, done, nil
}

func MustNewLazy() (resCPU, resRAM resources.ResourceViewer, done chan struct{}) {
	var err error
	resCPU, resRAM, done, err = NewLazy()
	if err != nil {
		panic(err)
	}
	return
}
