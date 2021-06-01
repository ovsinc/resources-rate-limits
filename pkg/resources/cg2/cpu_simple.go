package cg2

import (
	"os"

	"github.com/ovsinc/resources-rate-limits/internal/utils"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
)

type CPUCG2Simple struct{}

func NewCPUSimple() (rescommon.ResourceViewer, error) {
	cpu := &CPUCG2Simple{}
	return cpu, nil
}

func (cg *CPUCG2Simple) info() (total uint64, used uint64, err error) {
	ftotal, err := os.Open(rescommon.CGroup2CPULimitPath)
	if err != nil {
		return 0, 0, err
	}
	defer ftotal.Close()

	fused, err := os.Open(rescommon.CGroup2CPUUsagePath)
	if err != nil {
		return 0, 0, err
	}
	defer fused.Close()

	return getCPUInfo(ftotal, fused)
}

func (cg *CPUCG2Simple) Used() float64 {
	total, used, err := cg.info()
	if err != nil {
		rescommon.DbgErrCommon("CPUCG2Simple", err)
		return rescommon.FailValue
	}

	var p float64
	if used > 0 {
		p = utils.Percent(float64(used)/1000, float64(total))
	}
	rescommon.DbgInfCPU("CPUCG2Simple", 0, used, 0, total, p)

	return p
}
