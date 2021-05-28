package resources

import (
	"github.com/ovsinc/resources-rate-limits/pkg/resources/cg2"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
	"github.com/ovsinc/resources-rate-limits/pkg/resources/os"
)

type Resourcer rescommon.Resourcer

func AutoCPU() (Resourcer, error) {
	rt := Check()

	var res Resourcer
	switch rt {
	case ResourceType_OS:
		res, _ = os.NewCPUSimple()

	case ResourceType_CG2:
		res, _ = cg2.NewCPUSimple()

	default:
		return nil, ErrNoResourcer
	}

	return res, nil
}

func AutoRAM() (Resourcer, error) {
	rt := Check()

	var res Resourcer
	switch rt {
	case ResourceType_OS:
		res, _ = os.NewCPUSimple()

	case ResourceType_CG2:
		res, _ = cg2.NewCPUSimple()

	default:
		return nil, ErrNoResourcer
	}

	return res, nil
}
