// +build cg2
// not in docker with CGroups2

package cg2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCPUCG2Simple_Used(t *testing.T) {
	cpu, err := NewCPUSimple()
	assert.Nil(t, err)
	_ = cpu.Used()
}
