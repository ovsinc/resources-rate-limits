package os

import (
	"os"

	"github.com/ovsinc/resources-rate-limits/internal/utils"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
)

type MemOSSimple struct{}

func NewMemSimple() (rescommon.ResourceViewer, error) {
	return &MemOSSimple{}, nil
}

func (mem *MemOSSimple) info() (uint64, uint64, error) {
	f, err := os.Open(rescommon.RAMFilenameInfoProc)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	total, used, err := getMemInfo(f)
	if err != nil {
		return 0, 0, err
	}

	return total, used, nil
}

func (mem *MemOSSimple) Used() float64 {
	total, used, err := mem.info()
	if err != nil {
		return 0
	}

	return utils.Percent(float64(used), float64(total))
}
