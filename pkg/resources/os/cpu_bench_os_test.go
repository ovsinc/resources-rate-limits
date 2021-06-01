// +build os

package os

import (
	stdos "os"
	"testing"
	"time"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
	"go.uber.org/atomic"

	"github.com/stretchr/testify/require"
)

func BenchmarkCPUOSLazy_Used_Sys(b *testing.B) {
	f, err := stdos.Open(rescommon.CPUfilenameInfoProc)
	require.Nil(b, err)
	defer f.Close()

	done := make(chan struct{})
	defer close(done)

	cnf := rescommon.NewResourceConfig(rescommon.ResourceType_OS, rescommon.CPUfilenameInfoProc)
	require.Nil(b, cnf.Init())
	defer cnf.Stop()

	cpu, err := NewCPULazy(done, cnf, 600*time.Millisecond)
	require.Nil(b, err)

	time.Sleep(1300 * time.Millisecond)

	u := cpu.Used()
	require.Greater(b, u, float64(0))
	require.Less(b, u, float64(100))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cpu.Used()
	}
}

func BenchmarkCPUOSLazy_info_Sys(b *testing.B) {
	f, err := stdos.Open(rescommon.CPUfilenameInfoProc)
	require.Nil(b, err)
	defer f.Close()

	done := make(chan struct{})
	defer close(done)

	cpu := &CPUOSLazy{
		f:           f,
		utilization: &atomic.Float64{},
		done:        done,
		dur:         500 * time.Millisecond,
	}

	_, _, err = cpu.info()
	require.Nil(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = cpu.info()
	}
}

func BenchmarkCPUOSSimple_Used_Sys(b *testing.B) {
	cpu, err := NewCPUSimple()
	require.Nil(b, err)

	time.Sleep(2 * time.Second)

	u := cpu.Used()
	require.Greater(b, u, float64(0))
	require.Less(b, u, float64(100))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cpu.Used()
	}
}

func BenchmarkMemOSSimple_info_Sys(b *testing.B) {
	mi := &MemOSSimple{}

	_, _, err := mi.info()
	require.Nil(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = mi.info()
	}
}
