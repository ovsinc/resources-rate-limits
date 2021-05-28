package resources

import (
	pkgos "os"

	"github.com/ovsinc/errors"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
)

var ErrNoResourcer = errors.New("no valid resourcer")

func check(memlim, cpulim string) bool {
	if _, err := pkgos.Stat(memlim); err != nil {
		return false
	}
	if _, err := pkgos.Stat(cpulim); err != nil {
		return false
	}
	return true
}

type ResourceType uint16

const (
	ResourceType_UNKNOWN ResourceType = iota

	ResourceType_CG1
	ResourceType_CG2
	ResourceType_OS

	ResourceType_ENDS
)

func Check() (res ResourceType) {
	switch {
	case check(rescommon.CGroupCPULimitPath, rescommon.CGroupMemLimitPath):
		res = ResourceType_CG1

	case check(rescommon.CGroup2CPULimitPath, rescommon.CGroup2MemLimitPath):
		res = ResourceType_CG2

	case check(rescommon.CPUfilenameInfoProc, rescommon.RAMFilenameInfoProc):
		res = ResourceType_OS
	}

	return res
}
