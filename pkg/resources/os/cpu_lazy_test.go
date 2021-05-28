package os

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	data1 = `cpu  695636 4341 199895 27915902 7970 31786 16794 0 0 0
	cpu0 34184 339 10252 1296456 353 1077 1570 0 0 0
	cpu1 37919 380 10769 1294946 378 1185 693 0 0 0
	cpu2 38275 395 10732 1292207 451 1399 2039 0 0 0
	cpu3 38566 372 10888 1294215 423 1013 535 0 0 0
	cpu4 38375 321 10814 1294723 442 1081 430 0 0 0
	cpu5 30476 299 11449 1293194 401 6597 1461 0 0 0
	cpu6 40492 206 11235 1291149 407 1166 647 0 0 0
	cpu7 38739 399 11198 1294028 442 1016 379 0 0`
	data0 = ""
)

type rsMocStatic struct {
	data []byte
	done bool
}

func (r *rsMocStatic) Read(p []byte) (n int, err error) {
	if r.done {
		return 0, io.EOF
	}
	buf := bytes.NewBuffer(r.data)
	r.done = true
	return buf.Read(p)
}

func (r *rsMocStatic) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (r *rsMocStatic) Close() error { return nil }

func newBufferStatic(data []byte) io.ReadSeekCloser {
	return &rsMocStatic{data: data}
}

//

func BenchmarkCPUOSLazy_info_mock(b *testing.B) {
	cpu := &CPUOSLazy{
		f: newBufferStatic([]byte(data1)),
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

	cpu, err := NewCPULazy(done, cnf, 1000*time.Millisecond)
	require.Nil(t, err)
	defer cpu.Stop()

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
				f: newBufferStatic([]byte(data0)),
			},
			wantTotal: 0,
			wantUsed:  0,
			wantErr:   true,
		},
		{
			name: "< 7 args",
			fields: fields{
				f: newBufferStatic([]byte("cpu  695636 4341 199895 27915902")),
			},
			wantTotal: 0,
			wantUsed:  0,
			wantErr:   true,
		},
		{
			name: "normal",
			fields: fields{
				f: newBufferStatic([]byte(data1)),
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
