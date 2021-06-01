package cg1

import "testing"

func TestCPUCG1Simple_Used(t *testing.T) {
	cpu, _ := NewCPUSimple()
	_ = cpu.Used()
}
