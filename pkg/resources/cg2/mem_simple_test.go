package cg2

import "testing"

func TestMemCG2Simple_Used(t *testing.T) {
	mem, _ := NewMemSimple()
	_ = mem.Used()
}
