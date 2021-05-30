package resourcesratelimits

import (
	"time"

	"github.com/ovsinc/resources-rate-limits/pkg/resources/cg2"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
	"github.com/ovsinc/resources-rate-limits/pkg/resources/os"
)

type (
	Resourcer      rescommon.Resourcer
	ResourceViewer rescommon.ResourceViewer
)

func AutoCPUSimple(t rescommon.ResourceType) (ResourceViewer, error) {
	var res ResourceViewer
	switch t {
	case rescommon.ResourceType_OS:
		res, _ = os.NewCPUSimple()

	case rescommon.ResourceType_CG2:
		res, _ = cg2.NewCPUSimple()

	default:
		return nil, rescommon.ErrNoResourcer
	}

	return res, nil
}

func AutoRAMSimple(t rescommon.ResourceType) (ResourceViewer, error) {
	var res ResourceViewer
	switch t {
	case rescommon.ResourceType_OS:
		res, _ = os.NewCPUSimple()

	case rescommon.ResourceType_CG2:
		res, _ = cg2.NewCPUSimple()

	default:
		return nil, rescommon.ErrNoResourcer
	}

	return res, nil
}

func AutoLazyRAM(
	rt rescommon.ResourceConfiger,
	done chan struct{}, dur time.Duration,
) (res ResourceViewer, err error) {
	switch rt.Type() {
	case rescommon.ResourceType_OS:
		res, err = os.NewMemLazy(done, rt, dur)

	case rescommon.ResourceType_CG2:
		res, err = cg2.NewMemLazy(done, rt, dur)

	default:
		return nil, rescommon.ErrNoResourcer
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

func AutoLazyCPU(
	rt rescommon.ResourceConfiger,
	done chan struct{}, dur time.Duration,
) (res ResourceViewer, err error) {
	switch rt.Type() {
	case rescommon.ResourceType_OS:
		res, err = os.NewCPULazy(done, rt, dur)

	case rescommon.ResourceType_CG2:
		res, err = cg2.NewCPULazy(done, rt, dur)

	default:
		return nil, rescommon.ErrNoResourcer
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}
