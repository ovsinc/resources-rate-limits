// +build os

package os

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemOSSimple_Used_Sys(t *testing.T) {
	mi, _ := NewMemSimple()
	require.NotNil(t, mi)

	u := mi.Used()
	assert.Greater(t, u, float64(0))
	assert.Less(t, u, float64(100))
}
