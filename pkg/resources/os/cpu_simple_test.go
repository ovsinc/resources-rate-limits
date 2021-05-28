package os

import (
	"sync"
	"testing"
	"time"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkCPUOSSimple_info_mock(b *testing.B) {
	cpu := &CPUOSSimple{}

	_, _, err := cpu.info()
	require.Nil(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = cpu.info()
	}
}

//

func TestCPUOSSimple_Used_Sys(t *testing.T) {
	cpu, err := NewCPUSimple()
	require.Nil(t, err)

	used := cpu.Used()
	assert.Greater(t, used, float64(0))
	assert.Less(t, used, float64(100))
}

func TestCPUOSSimple_info(t *testing.T) {
	cpu := &CPUOSSimple{
		mu: new(sync.Mutex),
	}

	err := cpu.init()
	require.Nil(t, err)

	// подождем немного для стабилизации
	time.Sleep(rescommon.CPUSleep)

	total, used, err := cpu.info()
	assert.Nil(t, err)
	assert.Greater(t, total, uint64(0))
	assert.Greater(t, used, uint64(0))
}
