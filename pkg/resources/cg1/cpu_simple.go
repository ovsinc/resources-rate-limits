package cg1

import (
	"os"

	"github.com/ovsinc/resources-rate-limits/internal/utils"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
)

type CPUCG1Simple struct {
	limit, used string
}

func NewCPUSimple() (rescommon.ResourceViewer, error) {
	cpu := &CPUCG1Simple{
		limit: rescommon.CGroup2CPULimitPath,
		used:  rescommon.CGroup2CPUUsagePath,
	}
	return cpu, nil
}

func (cpu *CPUCG1Simple) info() (total uint64, used uint64, err error) {
	flimit, err := os.Open(cpu.limit)
	if err != nil {
		return 0, 0, err
	}
	defer flimit.Close()

	fusage, err := os.Open(cpu.used)
	if err != nil {
		return 0, 0, err
	}
	defer fusage.Close()

	total, err = readInfo(flimit)
	if err != nil {
		return 0, 0, err
	}

	used, err = readInfo(fusage)
	if err != nil {
		return 0, 0, err
	}

	return total, used, nil
}

func (cg *CPUCG1Simple) Used() float64 {
	total, used, err := cg.info()
	if err != nil {
		rescommon.DbgErrCommon("CPUCG1Simple", err)
		return rescommon.FailValue
	}

	var p float64
	if used > 0 {
		p = utils.Percent(float64(used)/1000000, float64(total))
	}
	rescommon.DbgInfCPU("CPUCG1Simple", 0, used, 0, total, p)

	return p
}
