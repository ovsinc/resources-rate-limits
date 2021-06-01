package cg1

import (
	"io"
	"testing"
	"time"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"

	"github.com/ovsinc/resources-rate-limits/pkg/resources/common/moc"
)

func BenchmarkCPUCG1Lazy_Used_moc(b *testing.B) {
	done := make(chan struct{})
	defer close(done)

	// cpuOK

	cpuOk := &CPUCG1Lazy{
		ftotal:      newMocStatic([]byte(CPUtotal)),
		fused:       newMocStatic([]byte(CPUused)),
		utilization: &atomic.Float64{},
		dur:         100 * time.Millisecond,
	}
	cpuOk.init()

	time.Sleep(300 * time.Millisecond)

	require.Equal(b, cpuOk.Used(), 0.053704441000000006)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cpuOk.Used()
	}
}

func BenchmarkCPUCG1Lazy_info_moc(b *testing.B) {
	cpu := &CPUCG1Lazy{
		ftotal: newMocStatic([]byte(CPUtotal)),
		fused:  newMocStatic([]byte(CPUused)),
	}

	total, used, err := cpu.info()
	require.Nil(b, err)
	require.Equal(b, total, uint64(100000))
	require.Equal(b, used, uint64(53704441))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = cpu.info()
	}
}

//

func TestCPUCG1Lazy_info_moc(t *testing.T) {
	type fields struct {
		dur         time.Duration
		ftotal      io.ReadSeekCloser
		fused       io.ReadSeekCloser
		utilization *atomic.Float64
		done        chan struct{}
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
				ftotal: newMocStatic([]byte("")),
				fused:  newMocStatic([]byte("")),
			},
			wantErr: true,
		},
		{
			name: "ok",
			fields: fields{
				ftotal: newMocStatic([]byte(CPUtotal)),
				fused:  newMocStatic([]byte(CPUused)),
			},
			wantTotal: 100000,
			wantUsed:  53704441,
		},
		{
			name: "empty used",
			fields: fields{
				ftotal: newMocStatic([]byte(CPUtotal)),
				fused:  newMocStatic([]byte("")),
			},
			wantErr: true,
		},
		{
			name: "empty total",
			fields: fields{
				fused:  newMocStatic([]byte(CPUused)),
				ftotal: newMocStatic([]byte("")),
			},
			wantErr: true,
		},
		{
			name: "unquoted",
			fields: fields{
				fused:  newMocStatic([]byte(CPUused)),
				ftotal: newMocStatic([]byte(CPUtotalUnquoted)),
			},
			wantUsed:  53704441,
			wantTotal: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cg := &CPUCG1Lazy{
				dur:         tt.fields.dur,
				ftotal:      tt.fields.ftotal,
				fused:       tt.fields.fused,
				utilization: tt.fields.utilization,
				done:        tt.fields.done,
			}
			gotTotal, gotUsed, err := cg.info()
			if (err != nil) != tt.wantErr {
				t.Errorf("CPUCG1Lazy.info() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("CPUCG1Lazy.info() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
			}
			if gotUsed != tt.wantUsed {
				t.Errorf("CPUCG1Lazy.info() gotUsed = %v, want %v", gotUsed, tt.wantUsed)
			}
		})
	}
}

func TestCPUCG1Lazy_Used_Moc(t *testing.T) {
	done := make(chan struct{})
	defer close(done)

	// cpuOK

	cpuOk := &CPUCG1Lazy{
		ftotal:      newMocStatic([]byte(CPUtotal)),
		fused:       newMocStatic([]byte(CPUused)),
		utilization: &atomic.Float64{},
		dur:         500 * time.Millisecond,
	}
	cpuOk.init()

	cpuFail := &CPUCG1Lazy{
		ftotal:      newMocStatic([]byte("")),
		fused:       newMocStatic([]byte("")),
		utilization: &atomic.Float64{},
		dur:         500 * time.Millisecond,
	}
	cpuFail.init()

	// wait

	time.Sleep(2 * time.Second)

	type fields struct {
		cpu *CPUCG1Lazy
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
			want: 0.053704441000000006,
		},
		{
			name: "fail",
			fields: fields{
				cpu: cpuFail,
			},
			want: float64(-1),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fields.cpu.Used()
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestNewCPULazy(t *testing.T) {
	done := make(chan struct{})
	defer close(done)

	cnf := moc.ResourceConfigMoc{
		Rtype: rescommon.ResourceType_CG1,
		FF:    make(map[string]io.ReadSeekCloser),
	}
	assert.Nil(t, cnf.Init())

	_, err := NewCPULazy(done, &cnf, 500*time.Millisecond)
	assert.NotNil(t, err)
}

func TestNewCPULazy_Moc(t *testing.T) {
	done := make(chan struct{})
	defer close(done)

	cnf := moc.ResourceConfigMoc{
		Rtype: rescommon.ResourceType_CG1,
		FF: map[string]io.ReadSeekCloser{
			rescommon.CGroupCPULimitPath: newMocStatic([]byte(CPUtotal)),
			rescommon.CGroupCPUUsagePath: newMocStatic([]byte(CPUused)),
		},
	}
	assert.Nil(t, cnf.Init())

	cpu, err := NewCPULazy(done, &cnf, 500*time.Millisecond)

	time.Sleep(time.Second)

	assert.Nil(t, err)
	assert.Equal(t, cpu.Used(), 0.053704441000000006)
}
