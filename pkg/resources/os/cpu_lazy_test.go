package os

import (
	"io"
	"os"
	"testing"
	"time"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkCPUOSLazy_info_mock(b *testing.B) {
	cpu := &CPUOSLazy{
		f: newMocStatic([]byte(data1)),
	}

	_, _, err := cpu.info()
	require.Nil(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = cpu.info()
	}
}

//

func TestNewCPULazy_Used_Sys(t *testing.T) {
	f, err := os.Open(rescommon.CPUfilenameInfoProc)
	require.Nil(t, err)
	defer f.Close()

	done := make(chan struct{})
	defer close(done)

	cnf := rescommon.NewResourceConfig(rescommon.ResourceType_OS, rescommon.CPUfilenameInfoProc)
	require.Nil(t, cnf.Init())
	defer cnf.Stop()

	cpu, err := NewCPULazy(done, cnf, 500*time.Millisecond)
	require.Nil(t, err)

	time.Sleep(1500 * time.Millisecond)

	used := cpu.Used()
	assert.Greater(t, used, float64(0))
	assert.Less(t, used, float64(100))
}

func TestCPUOSLazy_info_mock(t *testing.T) {
	type fields struct {
		f io.ReadSeekCloser
	}
	tests := []struct {
		name      string
		fields    fields
		wantTotal uint64
		wantUsed  uint64
		wantErr   bool
	}{
		{
			name: "empty",
			fields: fields{
				f: newMocStatic([]byte(data0)),
			},
			wantTotal: 0,
			wantUsed:  0,
			wantErr:   true,
		},
		{
			name: "< 7 args",
			fields: fields{
				f: newMocStatic([]byte("cpu  695636 4341 199895 27915902")),
			},
			wantTotal: 0,
			wantUsed:  0,
			wantErr:   true,
		},
		{
			name: "normal",
			fields: fields{
				f: newMocStatic([]byte(data1)),
			},
			wantTotal: 28872324,
			wantUsed:  948452,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cpu := &CPUOSLazy{
				f: tt.fields.f,
			}
			gotTotal, gotUsed, err := cpu.info()
			if (err != nil) != tt.wantErr {
				t.Errorf("CPUOSLazy.info() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("CPUOSLazy.info() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
			}
			if gotUsed != tt.wantUsed {
				t.Errorf("CPUOSLazy.info() gotUsed = %v, want %v", gotUsed, tt.wantUsed)
			}
		})
	}
}
