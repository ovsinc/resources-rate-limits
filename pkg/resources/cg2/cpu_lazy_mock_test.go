package cg2

import (
	"testing"
	"time"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
	resmoc "github.com/ovsinc/resources-rate-limits/pkg/resources/common/moc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
)

func BenchmarkCPUCG2Lazy_Used_moc(b *testing.B) {
	done := make(chan struct{})
	defer close(done)

	cpu := &CPUCG2Lazy{
		ftotal:      newCPUBufferStatic([]byte(cpuTotalMax)),
		fused:       newCPUBufferStatic([]byte(cpuStat)),
		utilization: &atomic.Float64{},
		dur:         100 * time.Millisecond,
		done:        done,
	}

	cpu.init()

	u := cpu.Used()
	require.Equal(b, u, 0.031585)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cpu.Used()
	}
}

func BenchmarkCPUCG2Lazy_info_moc(b *testing.B) {
	done := make(chan struct{})
	defer close(done)

	cpu := &CPUCG2Lazy{
		ftotal:      newCPUBufferStatic([]byte(cpuTotalMax)),
		fused:       newCPUBufferStatic([]byte(cpuStat)),
		utilization: &atomic.Float64{},
		dur:         100 * time.Millisecond,
		done:        done,
	}

	cpu.init()

	total, used, err := cpu.info()
	require.Nil(b, err)
	require.Equal(b, total, uint64(100000))
	require.Equal(b, used, uint64(31585))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = cpu.info()
	}
}

//

func TestCPUCG2Lazy_info_moc(t *testing.T) {
	type fields struct {
		ftotal rescommon.ReadSeekCloser
		fused  rescommon.ReadSeekCloser
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

func TestCPUCG2Lazy_Used_moc(t *testing.T) {
	done := make(chan struct{})
	cpuDone := &CPUCG2Lazy{
		ftotal:      newCPUBufferStatic([]byte(cpuTotalMax)),
		fused:       newCPUBufferStatic([]byte(cpuStat)),
		utilization: &atomic.Float64{},
		dur:         100 * time.Millisecond,
		done:        done,
	}
	cpuDone.init()
	close(done)

	cpuOk := &CPUCG2Lazy{
		ftotal:      newCPUBufferStatic([]byte(cpuTotalMax)),
		fused:       newCPUBufferStatic([]byte(cpuStat)),
		utilization: &atomic.Float64{},
		dur:         10 * time.Millisecond,
	}
	cpuOk.init()

	cpuFail := &CPUCG2Lazy{
		ftotal:      newCPUBufferStatic([]byte("")),
		fused:       newCPUBufferStatic([]byte("")),
		utilization: &atomic.Float64{},
		dur:         10 * time.Millisecond,
	}
	cpuFail.init()

	time.Sleep(15 * time.Millisecond)

	type fields struct {
		cpu *CPUCG2Lazy
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "ok",
			fields: fields{
				cpu: cpuOk,
			},
			want: 0.031585,
		},
		{
			name: "fail",
			fields: fields{
				cpu: cpuFail,
			},
			want: -1.0,
		},
		{
			name: "done",
			fields: fields{
				cpu: cpuDone,
			},
			want: -2.0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.cpu.Used(); got != tt.want {
				t.Errorf("CPUCG2Lazy.Used() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCPULazy(t *testing.T) {
	done := make(chan struct{})
	defer close(done)

	cnf := resmoc.ResourceConfigMoc{
		Rtype: rescommon.ResourceType_CG2,
		FF: map[string]rescommon.ReadSeekCloser{
			rescommon.CGroup2CPULimitPath: newCPUBufferStatic([]byte(cpuTotalMax)),
			rescommon.CGroup2CPUUsagePath: newCPUBufferStatic([]byte(cpuStat)),
		},
	}

	mem, err := NewCPULazy(done, &cnf, 100*time.Millisecond)
	assert.Nil(t, err)

	u := mem.Used()
	assert.Equal(t, u, 0.031585)

	_, err = NewCPULazy(done, nil, 100*time.Millisecond)
	assert.Error(t, err, rescommon.ErrNoResourceConfig)

	_, err = NewCPULazy(done, &cnf, 0)
	assert.Error(t, err, rescommon.ErrTickPeriodZero)

	cnf = resmoc.ResourceConfigMoc{
		Rtype: rescommon.ResourceType_CG2,
		FF:    map[string]rescommon.ReadSeekCloser{},
	}

	_, err = NewCPULazy(done, &cnf, 100*time.Millisecond)
	assert.NotNil(t, err)
}
