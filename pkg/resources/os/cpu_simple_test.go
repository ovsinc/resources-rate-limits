package os

import (
	"sync"
	"testing"
)

func TestNewCPUSimple(t *testing.T) {
	_, _ = NewCPUSimple()

	cpu := &CPUOSSimple{
		mu: new(sync.Mutex),
	}

	cpu.Used()
}
