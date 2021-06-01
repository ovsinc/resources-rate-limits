package cg1

import "testing"

func TestMemCG1Simple_Used(t *testing.T) {
	mem, _ := NewMemSimple()
	_ = mem.Used()
}
