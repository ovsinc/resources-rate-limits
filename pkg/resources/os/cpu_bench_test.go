package os_test

import (
	stdos "os"
	"testing"
	"time"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
	"github.com/ovsinc/resources-rate-limits/pkg/resources/os"

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

	cpu, err := os.NewCPULazy(done, cnf, 600*time.Millisecond)
	require.Nil(b, err)
	defer cpu.Stop()

	u := cpu.Used()
	require.Greater(b, u, float64(0))
	require.Less(b, u, float64(100))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cpu.Used()
	}
}

func BenchmarkCPUOSSimple_Used_Sys(b *testing.B) {
	cpu, err := os.NewCPUSimple()
	require.Nil(b, err)

	u := cpu.Used()
	require.Greater(b, u, float64(0))
	require.Less(b, u, float64(100))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cpu.Used()
	}
}
