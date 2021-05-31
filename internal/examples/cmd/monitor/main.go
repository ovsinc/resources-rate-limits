package main

import (
	"fmt"
	"time"

	ratelimits "github.com/ovsinc/resources-rate-limits"
)

func main() {
	cpu, ram, done := ratelimits.MustNewLazy()
	defer close(done)

	tick := time.NewTicker(500 * time.Millisecond)

	for range tick.C {
		fmt.Printf(
			"\rResources. CPU: %.2f. Ram: %.2f.",
			cpu.Used(), ram.Used(),
		)
	}
}
