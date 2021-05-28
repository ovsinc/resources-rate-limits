package os

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkMemOSSimple_info(b *testing.B) {
	mi := &MemOSSimple{}

	_, _, err := mi.info()
	require.Nil(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = mi.info()
	}
}

func TestMemOSSimple_Used_Sys(t *testing.T) {
	mi, _ := NewMemSimple()
	require.NotNil(t, mi)

	u := mi.Used()
	assert.Greater(t, u, float64(0))
	assert.Less(t, u, float64(100))
}

func TestMemOSSimple_info_mock(t *testing.T) {
	mem := &MemOSSimple{}

	total, used, err := mem.info()
	assert.Nil(t, err)
	assert.Greater(t, total, uint64(0))
	assert.Greater(t, used, uint64(0))
}
