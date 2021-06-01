package os

import (
	"testing"
)

func TestNewMemSimple(t *testing.T) {
	_, _ = NewMemSimple()

	mem := MemOSSimple{}

	_ = mem.Used()
}
