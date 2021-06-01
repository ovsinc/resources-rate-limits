package cg2

import (
	"testing"
)

func TestCPUCG2Simple_Used(t *testing.T) {
	cpu, _ := NewCPUSimple()
	_ = cpu.Used()
}
