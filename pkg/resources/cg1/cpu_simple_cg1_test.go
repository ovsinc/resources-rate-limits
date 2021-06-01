// +build cg1
// not in docker with CGroups

package cg1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCG1CPUCG1Simple_Used(t *testing.T) {
	cpu, err := NewCPUSimple()
	require.Nil(t, err)

	u := cpu.Used()
	assert.Greater(t, u, float64(0))
	assert.Less(t, u, float64(100))
}
