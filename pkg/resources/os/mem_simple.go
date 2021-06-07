package os

import (
	"os"

	"github.com/ovsinc/resources-rate-limits/internal/utils"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
)

type MemOSSimple struct {
	stat string
}

func NewMemSimple() (rescommon.ResourceViewer, error) {
	return &MemOSSimple{
		stat: rescommon.RAMFilenameInfoProc,
	}, nil
}

func (mem *MemOSSimple) info() (uint64, uint64, error) {
	f, err := os.Open(mem.stat)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	total, used, err := GetMemInfo(f)
	if err != nil {
		return 0, 0, err
	}

	return total, used, nil
}

func (mem *MemOSSimple) Used() float64 {
	total, used, err := mem.info()
	if err != nil {
		rescommon.DbgErrCommon("MemOSSimple", err)
		return rescommon.FailValue
	}

	p := utils.Percent(float64(used), float64(total))
	rescommon.DbgInfRAM("MemOSSimple", used, total, p)

	return p
}
