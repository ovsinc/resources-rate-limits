package os

import (
	"os"
	"sync"
	"time"

	"github.com/ovsinc/resources-rate-limits/internal/utils"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
)

var _ rescommon.Resourcer = (*CPUOSSimple)(nil)

type CPUOSSimple struct {
	mu                  *sync.Mutex
	prevTotal, prevUsed uint64
}

func NewCPUSimple() (*CPUOSSimple, error) {
	cpu := &CPUOSSimple{
		mu: new(sync.Mutex),
	}

	err := cpu.init()
	if err != nil {
		return nil, err
	}

	// подождем немного для стабилизации
	time.Sleep(rescommon.CPUSleep)

	return cpu, nil
}

func (cpu *CPUOSSimple) Used() float64 {
	total, used, err := cpu.info()
	if err != nil {
		return 0
	}

	cpu.mu.Lock()
	defer cpu.mu.Unlock()

	percent := utils.CPUPercent(cpu.prevUsed, used, cpu.prevTotal, total)

	cpu.prevUsed = used
	cpu.prevTotal = total

	return percent
}

func (cpu *CPUOSSimple) Stop() {}

func (cpu *CPUOSSimple) info() (total uint64, used uint64, err error) {
	f, err := os.Open(rescommon.CPUfilenameInfoProc)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	return getCPUInfo(f)
}

func (cpu *CPUOSSimple) init() error {
	total, used, err := cpu.info()
	if err != nil {
		return err
	}

	cpu.prevUsed = used
	cpu.prevTotal = total

	return nil
}
