package cg1

import (
	"os"

	"github.com/ovsinc/resources-rate-limits/internal/utils"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
)

type MemCG1Simple struct{}

func NewMemSimple() (rescommon.ResourceViewer, error) {
	return &MemCG1Simple{}, nil
}

func (mem *MemCG1Simple) info() (uint64, uint64, error) {
	flimit, err := os.Open(rescommon.CGroup2MemLimitPath)
	if err != nil {
		return 0, 0, err
	}
	defer flimit.Close()

	fused, err := os.Open(rescommon.CGroup2MemUsagePath)
	if err != nil {
		return 0, 0, err
	}
	defer fused.Close()

	limit, err := readInfo(flimit)
	if err != nil {
		return 0, 0, err
	}

	used, err := readInfo(flimit)
	if err != nil {
		return 0, 0, err
	}

	return limit, used, nil
}

func (mem *MemCG1Simple) Stop() {}

func (mem *MemCG1Simple) Used() float64 {
	total, used, err := mem.info()
	if err != nil {
		return rescommon.FailValue
	}

	return utils.Percent(float64(used), float64(total))
}
