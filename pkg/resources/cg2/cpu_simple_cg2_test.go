// +build cg2
// not in docker with CGroups2

package cg2_test

import (
	"testing"

	"github.com/ovsinc/resources-rate-limits/pkg/resources/cg2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCG2CPUCG2Simple_Used(t *testing.T) {
	cpu, err := cg2.NewCPUSimple()
	require.Nil(t, err)

	u := cpu.Used()
	assert.Greater(t, u, float64(0))
	assert.Less(t, u, float64(100))
}
