package resourcesratelimits

import (
	pkgos "os"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
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
