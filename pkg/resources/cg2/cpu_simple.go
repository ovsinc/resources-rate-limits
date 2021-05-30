package cg2

import (
	"os"
	"sync"
	"time"

	"github.com/ovsinc/resources-rate-limits/internal/utils"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
)

type CPUCG2Simple struct {
	mu                  *sync.Mutex
	prevTotal, prevUsed uint64
}

func NewCPUSimple() (rescommon.ResourceViewer, error) {
	cpu := &CPUCG2Simple{
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

func (cg *CPUCG2Simple) init() error {
	total, used, err := cg.info()
	if err != nil {
		return err
	}

	cg.prevUsed = used
	cg.prevTotal = total

	return nil
}

func (cg *CPUCG2Simple) Used() float64 {
	total, used, err := cg.info()
	if err != nil {
		rescommon.Debug("[CPUCG2Simple]<ERR> Check resource fails with %v", err)
		return rescommon.FailValue
	}

	cg.mu.Lock()
	defer cg.mu.Unlock()

	percent := utils.CPUPercent(cg.prevUsed, used, cg.prevTotal, total)

	cg.prevUsed = used
	cg.prevTotal = total

	rescommon.Debug(
		"[CPUCG2Simple]<INFO> last: %d/%d now: %d/%d",
		cg.prevUsed, cg.prevTotal, used, total,
	)

	return percent
}
