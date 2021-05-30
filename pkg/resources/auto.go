package resources

import (
	pkgos "os"
	"time"

	"github.com/ovsinc/resources-rate-limits/pkg/resources/cg2"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
	"github.com/ovsinc/resources-rate-limits/pkg/resources/os"
)

type (
	Resourcer      rescommon.Resourcer
	ResourceViewer rescommon.ResourceViewer
)

func check(files ...string) bool {
	for _, file := range files {
		if _, err := pkgos.Stat(file); err != nil {
			return false
		}
	}
	return true
}

func Check() rescommon.ResourceConfiger {
	var (
		files = []string{}
		t     rescommon.ResourceType
	)

	switch {
	case check(rescommon.CGroupCPULimitPath, rescommon.CGroupMemLimitPath):
		t = rescommon.ResourceType_CG1
		files = rescommon.CGroupFiles

	case check(rescommon.CGroup2CPULimitPath, rescommon.CGroup2MemLimitPath):
		t = rescommon.ResourceType_CG2
		files = rescommon.CGroup2Files

	case check(rescommon.CPUfilenameInfoProc, rescommon.RAMFilenameInfoProc):
		t = rescommon.ResourceType_OS
		files = rescommon.OSLinuxFiles
	}

	return rescommon.NewResourceConfig(t, files...)
}

func AutoCPUSimple() (ResourceViewer, error) {
	rt := Check()

	var res ResourceViewer
	switch rt.Type() {
	case rescommon.ResourceType_OS:
		res, _ = os.NewCPUSimple()

	case rescommon.ResourceType_CG2:
		res, _ = cg2.NewCPUSimple()

	default:
		return nil, rescommon.ErrNoResourcer
	}

	return res, nil
}

func AutoRAMSimple() (ResourceViewer, error) {
	rt := Check()

	var res ResourceViewer
	switch rt.Type() {
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
	done chan struct{}, dur time.Duration,
) (res ResourceViewer, err error) {
	rt := Check()

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
	done chan struct{}, dur time.Duration,
) (res ResourceViewer, err error) {
	rt := Check()

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
