package common

import (
	"time"

	"github.com/ovsinc/errors"
)

const (
	RAMFilenameInfoProc = "/proc/meminfo"
	CPUfilenameInfoProc = "/proc/stat"

	// https://www.kernel.org/doc/html/latest/admin-guide/cgroup-v2.html
	CGroup2MemUsagePath = "/sys/fs/cgroup/memory.current"
	CGroup2MemLimitPath = "/sys/fs/cgroup/memory.max"

	CGroup2CPUUsagePath = "/sys/fs/cgroup/cpu.stat"
	CGroup2CPULimitPath = "/sys/fs/cgroup/cpu.max"

	CGroupMemLimitPath = "/sys/fs/cgroup/memory/memory.limit_in_bytes"
	CGroupMemUsagePath = "/sys/fs/cgroup/memory/memory.usage_in_bytes"

	CGroupCPULimitPath = "/sys/fs/cgroup/cpu/cpu.cfs_quota_us"
	CGroupCPUUsagePath = "/sys/fs/cgroup/cpu/cpuacct.usage"

	CPUSleep = time.Second

	DefaultDuration = 3 * time.Second

	FailValue = -1.0
	DoneValue = -2.0
)

var (
	ErrAllIsZero          = errors.New("total is zero")
	ErrTickPeriodZero     = errors.New("check period is zero")
	ErrNoResourceConfig   = errors.New("no resource config")
	ErrNoResourceReadFile = errors.New("file for check resource info not found")
)
