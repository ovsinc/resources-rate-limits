package cg2

import (
	"io"
	"testing"
	"time"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
	resmoc "github.com/ovsinc/resources-rate-limits/pkg/resources/common/moc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
)

func BenchmarkCPUCG2Lazy_info_mock(b *testing.B) {
	cpu := &CPUCG2Lazy{
		ftotal: newCPUBufferStatic([]byte(cpuTotalMax)),
		fused:  newCPUBufferStatic([]byte(cpuStat)),
	}

	_, _, err := cpu.info()
	require.Nil(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = cpu.info()
	}
}

//

func TestCPUCG2Lazy_info_mock(t *testing.T) {
	type fields struct {
		ftotal io.ReadSeekCloser
		fused  io.ReadSeekCloser
	}
	tests := []struct {
		name      string
		fields    fields
		wantTotal uint64
		wantUsed  uint64
		wantErr   bool
	}{
		{
			name: "ok",
			fields: fields{
				ftotal: newCPUBufferStatic([]byte(cpuTotalMax)),
				fused:  newCPUBufferStatic([]byte(cpuStat)),
			},
			wantTotal: 100000,
			wantUsed:  31585,
		},
		{
			name: "empty",
			fields: fields{
				ftotal: newCPUBufferStatic([]byte("")),
				fused:  newCPUBufferStatic([]byte("")),
			},
			wantErr: true,
		},
		{
			name: "empty used",
			fields: fields{
				ftotal: newCPUBufferStatic([]byte(cpuTotalMax)),
				fused:  newCPUBufferStatic([]byte("")),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cg := &CPUCG2Lazy{
				ftotal: tt.fields.ftotal,
				fused:  tt.fields.fused,
			}
			gotTotal, gotUsed, err := cg.info()
			if (err != nil) != tt.wantErr {
				t.Errorf("Cg2cpu.Info() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("Cg2cpu.Info() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
			}
			if gotUsed != tt.wantUsed {
				t.Errorf("Cg2cpu.Info() gotUsed = %v, want %v", gotUsed, tt.wantUsed)
			}
		})
	}
}

func TestCPUCG2Lazy_Used_Moc(t *testing.T) {
	cpu := &CPUCG2Lazy{
		ftotal:      newCPUBufferStatic([]byte(cpuTotalMax)),
		fused:       newCPUBufferStatic([]byte(cpuStat)),
		utilization: &atomic.Float64{},
		dur:         500 * time.Millisecond,
	}

	done := make(chan struct{})
	defer close(done)

	cpu.init()

	time.Sleep(2 * time.Second)

	u := cpu.Used()
	assert.Equal(t, u, float64(0))
}

func TestNewCPULazy(t *testing.T) {
	done := make(chan struct{})
	defer close(done)

	cnf := &resmoc.ResourceConfigMoc{
		Rtype: rescommon.ResourceType_CG2,
		FF: map[string]io.ReadSeekCloser{
			rescommon.CGroup2CPULimitPath: newCPUBufferStatic([]byte(cpuTotalMax)),
			rescommon.CGroup2CPUUsagePath: newCPUBufferStatic([]byte(cpuStat)),
		},
	}
	assert.Nil(t, cnf.Init())

	mem, err := NewCPULazy(done, cnf, 100*time.Millisecond)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	u := mem.Used()
	assert.Equal(t, u, float64(0))
}
