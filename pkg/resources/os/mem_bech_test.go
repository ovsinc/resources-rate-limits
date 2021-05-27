package os_test

import (
	stdos "os"
	"testing"
	"time"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
	"github.com/ovsinc/resources-rate-limits/pkg/resources/os"

	"github.com/stretchr/testify/require"
)

func BenchmarkMemOSLazy_Used_Sys(b *testing.B) {
	f, err := stdos.Open(rescommon.RAMFilenameInfoProc)
	require.Nil(b, err)
	defer f.Close()

	done := make(chan struct{})
	defer close(done)

	mem, err := os.NewMemLazy(done, f, 10*time.Millisecond)
	require.Nil(b, err)
	defer mem.Stop()

	u := mem.Used()
	require.Greater(b, u, float64(0))
	require.Less(b, u, float64(100))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mem.Used()
	}
}

func BenchmarkMemOSSimple_Used_Sys(b *testing.B) {
	mem := os.NewMemSimple()

	u := mem.Used()
	require.Greater(b, u, float64(0))
	require.Less(b, u, float64(100))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mem.Used()
	}
}
