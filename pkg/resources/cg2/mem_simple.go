package cg2

import (
	"os"

	"github.com/ovsinc/resources-rate-limits/internal/utils"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
)

type MemCG2Simple struct {
	total, used string
	proc        string
}

func NewMemSimple() (rescommon.ResourceViewer, error) {
	return &MemCG2Simple{
		total: rescommon.CGroup2MemLimitPath,
		used:  rescommon.CGroup2CPUUsagePath,
		proc:  rescommon.RAMFilenameInfoProc,
	}, nil
}

func (cg *MemCG2Simple) info() (uint64, uint64, error) {
	ftotal, err := os.Open(cg.total)
	if err != nil {
		return 0, 0, err
	}
	defer ftotal.Close()

	fused, err := os.Open(cg.used)
	if err != nil {
		return 0, 0, err
	}
	defer fused.Close()

	fprocmem, err := os.Open(cg.proc)
	if err != nil {
		return 0, 0, err
	}
	defer fprocmem.Close()

	return getMemInfo(ftotal, fused, fprocmem)
}

func (cg *MemCG2Simple) Used() float64 {
	total, used, err := cg.info()
	if err != nil {
		rescommon.DbgErrCommon("MemCG2Simple", err)
		return rescommon.FailValue
	}

	p := utils.Percent(float64(used), float64(total))
	rescommon.DbgInfRAM("MemCG2Simple", used, total, p)

	return p
}
