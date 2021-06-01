package os

import (
	"testing"
	"time"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
	resmoc "github.com/ovsinc/resources-rate-limits/pkg/resources/common/moc"
	"go.uber.org/atomic"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkCPUOSLazy_info_moc(b *testing.B) {
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

func TestCPUOSLazy_info_moc(t *testing.T) {
	type fields struct {
		f rescommon.ReadSeekCloser
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
			wantTotal: 13133069,
			wantUsed:  734265,
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

func TestCPUOSLazy_Used_moc(t *testing.T) {
	cpuOk := &CPUOSLazy{
		f:           newMocStatic([]byte(data1)),
		utilization: &atomic.Float64{},
		dur:         10 * time.Millisecond,
	}
	cpuOk.init()

	cpuFail := &CPUOSLazy{
		f:           newMocStatic([]byte("")),
		utilization: &atomic.Float64{},
		dur:         10 * time.Millisecond,
	}
	cpuFail.init()

	done := make(chan struct{})
	cpuDone := &CPUOSLazy{
		f:           newMocStatic([]byte(data1)),
		utilization: &atomic.Float64{},
		dur:         100 * time.Millisecond,
		done:        done,
	}
	cpuDone.init()
	close(done)

	time.Sleep(20 * time.Millisecond)

	type fields struct {
		cpu *CPUOSLazy
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
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.cpu.Used(); got != tt.want {
				t.Errorf("CPUOSLazy.Used() = %v, want %v", got, tt.want)
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
			rescommon.CPUfilenameInfoProc: newMocStatic([]byte(data1)),
		},
	}

	cpu, err := NewCPULazy(done, &cnf, 100*time.Millisecond)
	assert.Nil(t, err)
	assert.Equal(t, cpu.Used(), float64(0))

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
