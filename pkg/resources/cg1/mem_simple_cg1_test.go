// +build cg1
// not in docker with CGroups

package cg1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCG1MemCG1Simple_Used(t *testing.T) {
	mem, err := NewMemSimple()
	require.Nil(t, err)

	u := mem.Used()
	assert.Greater(t, u, float64(0))
	assert.Less(t, u, float64(100))
}
