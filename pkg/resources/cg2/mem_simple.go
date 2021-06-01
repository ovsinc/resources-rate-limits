package cg2

import (
	"os"

	"github.com/ovsinc/resources-rate-limits/internal/utils"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
)

type MemCG2Simple struct{}

func NewMemSimple() (rescommon.ResourceViewer, error) {
	return &MemCG2Simple{}, nil
}

func (cg *MemCG2Simple) info() (uint64, uint64, error) {
	ftotal, err := os.Open(rescommon.CGroup2MemLimitPath)
	if err != nil {
		return 0, 0, err
	}
	defer ftotal.Close()

	fused, err := os.Open(rescommon.CGroup2CPUUsagePath)
	if err != nil {
		return 0, 0, err
	}
	defer fused.Close()

	fprocmem, err := os.Open(rescommon.RAMFilenameInfoProc)
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
