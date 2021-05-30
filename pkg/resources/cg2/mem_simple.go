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

	return getMemInfo(ftotal, fused)
}

func (cg *MemCG2Simple) Used() float64 {
	total, used, err := cg.info()
	if err != nil {
		rescommon.Debug("[MemCG2Simple]<ERR> Check resource fails with %v", err)
		return rescommon.FailValue
	}

	rescommon.Debug(
		"[MemCG2Simple]<INFO> now: %d/%d",
		used, total,
	)

	return utils.Percent(float64(used), float64(total))
}
