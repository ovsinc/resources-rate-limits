package cg1

import (
	"os"
	"sync"
	"time"

	"github.com/ovsinc/resources-rate-limits/internal/utils"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
)

type CPUCG1Simple struct {
	mu                  *sync.Mutex
	prevTotal, prevUsed uint64
}

func NewCPUSimple() (rescommon.ResourceViewer, error) {
	cpu := &CPUCG1Simple{
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

func (cg *CPUCG1Simple) init() error {
	total, used, err := cg.info()
	if err != nil {
		return err
	}

	cg.prevUsed = used
	cg.prevTotal = total

	return nil
}

func (cpu *CPUCG1Simple) info() (total uint64, used uint64, err error) {
	flimit, err := os.Open(rescommon.CGroup2CPULimitPath)
	if err != nil {
		return 0, 0, err
	}
	defer flimit.Close()

	fusage, err := os.Open(rescommon.CGroup2CPUUsagePath)
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

func (cg *CPUCG1Simple) Stop() {}

func (cg *CPUCG1Simple) Used() float64 {
	total, used, err := cg.info()
	if err != nil {
		return 0
	}

	cg.mu.Lock()
	defer cg.mu.Unlock()

	percent := utils.CPUPercent(cg.prevUsed, used, cg.prevTotal, total)

	cg.prevUsed = used
	cg.prevTotal = total

	return percent
}
