// +build cg2
// not in docker with CGroups2

package cg2_test

import (
	"testing"

	"github.com/ovsinc/resources-rate-limits/pkg/resources/cg2"

	"github.com/stretchr/testify/assert"
)

func TestCG2MemCG2Simple_Used(t *testing.T) {
	mem, err := cg2.NewMemSimple()
	assert.Nil(t, err)

	u := mem.Used()
	assert.Greater(t, u, float64(0))
	assert.Less(t, u, float64(100))
}
