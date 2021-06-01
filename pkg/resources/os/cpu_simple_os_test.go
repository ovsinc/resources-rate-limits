// +build os

package os

import (
	"sync"
	"testing"
	"time"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//

func TestCPUOSSimple_Used_Sys(t *testing.T) {
	cpu, err := NewCPUSimple()
	require.Nil(t, err)

	time.Sleep(2 * rescommon.CPUSleep)

	used := cpu.Used()
	assert.Greater(t, used, float64(0))
	assert.Less(t, used, float64(100))
}

func TestCPUOSSimple_info_Sys(t *testing.T) {
	cpu := &CPUOSSimple{
		mu: new(sync.Mutex),
	}

	err := cpu.init()
	require.Nil(t, err)

	time.Sleep(2 * rescommon.CPUSleep)

	total, used, err := cpu.info()
	assert.Nil(t, err)
	assert.Greater(t, total, uint64(0))
	assert.Greater(t, used, uint64(0))
}
