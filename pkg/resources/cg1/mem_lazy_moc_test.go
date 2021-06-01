package cg1

import (
	"io"
	"testing"
	"time"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
	"github.com/ovsinc/resources-rate-limits/pkg/resources/common/moc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
)

func BenchmarkMemCG1Lazy_Used_moc(b *testing.B) {
	done := make(chan struct{})
	defer close(done)

	mem := &MemCG1Lazy{
		ftotal: newMocStatic([]byte(MemTotal)),
		fused:  newMocStatic([]byte(MemUsed)),
		dur:    100 * time.Millisecond,
		used:   &atomic.Float64{},
		done:   done,
	}
	mem.init()

	time.Sleep(300 * time.Millisecond)

	u := mem.Used()
	require.Equal(b, u, 2.8125)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mem.Used()
	}
}

func BenchmarkMemCG1Lazy_info_moc(b *testing.B) {
	mem := &MemCG1Lazy{
		ftotal: newMocStatic([]byte(MemTotal)),
		fused:  newMocStatic([]byte(MemUsed)),
	}

	total, used, err := mem.info()
	require.Nil(b, err)
	require.Equal(b, total, uint64(10485760))
	require.Equal(b, used, uint64(294912))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = mem.info()
	}
}

//

func TestMemCG1Lazy_info_moc(t *testing.T) {
	type fields struct {
		ftotal io.ReadSeekCloser
		fused  io.ReadSeekCloser
		used   *atomic.Float64
		dur    time.Duration
		done   chan struct{}
	}
	tests := []struct {
		name      string
		fields    fields
		wantTotal uint64
		wantUsed  uint64
		wantErr   bool
	}{
		{
			name: "normal",
			fields: fields{
				ftotal: newMocStatic([]byte(MemTotal)),
				fused:  newMocStatic([]byte(MemUsed)),
			},
			wantTotal: 10485760,
			wantUsed:  294912,
		},
		{
			name: "fails total",
			fields: fields{
				ftotal: newMocStatic([]byte(MemTotalFail)),
				fused:  newMocStatic([]byte(MemUsed)),
			},
			wantErr: true,
		},
		{
			name: "fails used",
			fields: fields{
				ftotal: newMocStatic([]byte(MemTotal)),
				fused:  newMocStatic([]byte(MemUsedFail)),
			},
			wantErr: true,
		},
		{
			name: "unqouted",
			fields: fields{
				ftotal: newMocStatic([]byte(MemTotalUnquoted)),
				fused:  newMocStatic([]byte(MemUsed)),
			},
			wantUsed:  294912,
			wantTotal: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cg := &MemCG1Lazy{
				ftotal: tt.fields.ftotal,
				fused:  tt.fields.fused,
				used:   tt.fields.used,
				dur:    tt.fields.dur,
				done:   tt.fields.done,
			}
			gotTotal, gotUsed, err := cg.info()
			if (err != nil) != tt.wantErr {
				t.Errorf("MemCG1Lazy.info() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("MemCG1Lazy.info() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
			}
			if gotUsed != tt.wantUsed {
				t.Errorf("MemCG1Lazy.info() gotUsed = %v, want %v", gotUsed, tt.wantUsed)
			}
		})
	}
}

func TestMemCG1Lazy_Used_moc(t *testing.T) {
	done := make(chan struct{})
	defer close(done)

	mem := &MemCG1Lazy{
		ftotal: newMocStatic([]byte(MemTotal)),
		fused:  newMocStatic([]byte(MemUsed)),
		dur:    100 * time.Millisecond,
		used:   &atomic.Float64{},
		done:   done,
	}
	mem.init()

	time.Sleep(300 * time.Millisecond)

	u := mem.Used()
	assert.Equal(t, u, 2.8125)
}

func TestNewMemLazy(t *testing.T) {
	done := make(chan struct{})
	defer close(done)

	_, err := NewMemLazy(done, nil, 100*time.Millisecond)
	assert.Error(t, err, rescommon.ErrNoResourceConfig)

	cnf := moc.ResourceConfigMoc{
		Rtype: rescommon.ResourceType_CG1,
		FF:    make(map[string]io.ReadSeekCloser),
	}

	_, err = NewMemLazy(done, &cnf, 0)
	assert.Error(t, err, rescommon.ErrTickPeriodZero)

	_, err = NewMemLazy(done, &cnf, 100*time.Millisecond)
	assert.NotNil(t, err)

	cnf = moc.ResourceConfigMoc{
		Rtype: rescommon.ResourceType_CG1,
		FF: map[string]io.ReadSeekCloser{
			rescommon.CGroupMemUsagePath: newMocStatic([]byte(MemUsed)),
			rescommon.CGroupMemLimitPath: newMocStatic([]byte(MemTotal)),
		},
	}

	_, err = NewMemLazy(done, &cnf, 100*time.Millisecond)
	assert.Nil(t, err)
}
