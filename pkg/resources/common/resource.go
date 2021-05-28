package common

import (
	"github.com/ovsinc/errors"
)

var ErrNoResourcer = errors.New("no valid resourcer")

type ResourceType uint16

const (
	ResourceType_UNKNOWN ResourceType = iota

	ResourceType_CG1
	ResourceType_CG2
	ResourceType_OS

	ResourceType_ENDS
)

var (
	CGroupFiles = []string{
		CGroupCPULimitPath,
		CGroupMemLimitPath,
		CGroupCPUUsagePath,
		CGroupMemUsagePath,
	}

	CGroup2Files = []string{
		CGroup2CPULimitPath,
		CGroup2MemLimitPath,
		CGroup2CPUUsagePath,
		CGroup2MemUsagePath,
	}

	OSLinuxFiles = []string{
		CPUfilenameInfoProc,
		RAMFilenameInfoProc,
	}
)
