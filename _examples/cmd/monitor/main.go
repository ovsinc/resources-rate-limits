package main

import (
	"fmt"
	"time"

	resos "github.com/ovsinc/resources-rate-limits/pkg/resources/os"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
	"github.com/ovsinc/resources-rate-limits/pkg/resources"
)

func main() {
	rescommon.
	cpu, _ := resos.NewCPULazy()
	defer cpu.Stop()

	ram, _ := resos.AutoRAM()
	defer ram.Stop()

	tick := time.NewTicker(500 * time.Millisecond)

	for range tick.C {
		fmt.Printf(
			"\rResources. CPU: %.2f. Ram: %.2f.",
			cpu.Used(), ram.Used(),
		)
	}
}
