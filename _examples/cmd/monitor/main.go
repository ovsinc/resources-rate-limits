package main

import (
	"fmt"
	"time"

	"github.com/ovsinc/resources-rate-limits/resources-rate-limits/internal/resources"
)

func main() {
	r := resources.MustNew(3*time.Second, 4*time.Second)
	defer r.Shutdown()

	tick := time.NewTicker(500 * time.Millisecond)

	for range tick.C {
		fmt.Printf(
			"\rResources. CPU: %.2f. Ram: %.2f.",
			r.CPUUtilization(), r.RAMUtilization(),
		)
	}
}
