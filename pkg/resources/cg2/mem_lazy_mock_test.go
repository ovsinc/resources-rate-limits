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

func BenchmarkMemCG2Lazy_info_moc(b *testing.B) {
	mem := &MemCG2Lazy{
		ftotal: newMemBufferStatic([]byte(memTotalData1)),
		fused:  newMemBufferStatic([]byte(memUsedData1)),
		dur:    100 * time.Millisecond,
	}

	total, used, err := mem.info()
	require.Nil(b, err)
	require.Equal(b, total, uint64(104857600))
	require.Equal(b, used, uint64(1482752))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = mem.info()
	}
}

func BenchmarkMemCG2Lazy_Used_moc(b *testing.B) {
	done := make(chan struct{})
	defer close(done)

	mem := &MemCG2Lazy{
		ftotal: newMemBufferStatic([]byte(memTotalData1)),
		fused:  newMemBufferStatic([]byte(memUsedData1)),
		used:   &atomic.Float64{},
		dur:    100 * time.Millisecond,
		done:   done,
	}

	mem.init()

	u := mem.Used()
	require.Equal(b, u, 1.4140625)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mem.Used()
	}
}

//

func TestNewMemLazy_mock(t *testing.T) {
	done := make(chan struct{})
	defer close(done)

	ftotal := newMemBufferStatic([]byte(memTotalData1))
	fused := newMemBufferStatic([]byte(memUsedData1))
	fprocmem := newMemBufferStatic([]byte(meminfo))

	cnf := resmoc.ResourceConfigMoc{
		Rtype: rescommon.ResourceType_CG2,
		FF: map[string]rescommon.ReadSeekCloser{
			rescommon.CGroup2MemLimitPath: ftotal,
			rescommon.CGroup2MemUsagePath: fused,
			rescommon.RAMFilenameInfoProc: fprocmem,
		},
	}

	_, err := NewMemLazy(done, &cnf, time.Millisecond)
	assert.Nil(t, err)

	_, err = NewMemLazy(done, &cnf, 0)
	assert.Error(t, err, rescommon.ErrTickPeriodZero)

	_, err = NewMemLazy(done, nil, time.Millisecond)
	assert.Error(t, err, rescommon.ErrNoResourceConfig)

	cnf = resmoc.ResourceConfigMoc{
		Rtype: rescommon.ResourceType_CG2,
		FF:    map[string]rescommon.ReadSeekCloser{},
	}

	_, err = NewMemLazy(done, &cnf, time.Millisecond)
	assert.NotNil(t, err)
}

func TestMemCG2Lazy_info_moc(t *testing.T) {
	type fields struct {
		ftotal   rescommon.ReadSeekCloser
		fused    rescommon.ReadSeekCloser
		fprocmem rescommon.ReadSeekCloser
	}
	tests := []struct {
		name      string
		fields    fields
		wantUsed  uint64
		wantTotal uint64
		wantErr   bool
	}{
		{
			name: "used zero",
			fields: fields{
				ftotal:   newMemBufferStatic([]byte(memTotalDataMax)),
				fused:    newMemBufferStatic([]byte(memUsedZero)),
				fprocmem: newMemBufferStatic([]byte(meminfo)),
			},
			wantErr: true,
		},
		{
			name: "total zero",
			fields: fields{
				ftotal:   newMemBufferStatic([]byte("")),
				fused:    newMemBufferStatic([]byte(memUsedData1)),
				fprocmem: newMemBufferStatic([]byte("")),
			},
			wantErr: true,
		},
		{
			name: "normal",
			fields: fields{
				ftotal:   newMemBufferStatic([]byte(memTotalData1)),
				fused:    newMemBufferStatic([]byte(memUsedData1)),
				fprocmem: newMemBufferStatic([]byte(meminfo)),
			},
			wantTotal: 104857600,
			wantUsed:  1482752,
		},
		{
			name: "max total",
			fields: fields{
				ftotal:   newMemBufferStatic([]byte(memTotalDataMax)),
				fused:    newMemBufferStatic([]byte(memUsedData1)),
				fprocmem: newMemBufferStatic([]byte(meminfo)),
			},
			wantTotal: 33114734592,
			wantUsed:  1482752,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cg := &MemCG2Lazy{
				ftotal:   tt.fields.ftotal,
				fused:    tt.fields.fused,
				fprocmem: tt.fields.fprocmem,
			}
			total, used, err := cg.info()
			if (err != nil) != tt.wantErr {
				t.Errorf("cg2mem.getMemInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if total != tt.wantTotal || used != tt.wantUsed {
				t.Errorf("cg2mem.getMemInfo() = %d | %d, want %d | %d", used, tt.wantUsed, total, tt.wantTotal)
			}
		})
	}
}

func TestMemCG2Lazy_Used_moc(t *testing.T) {
	done := make(chan struct{})
	memDone := &MemCG2Lazy{
		ftotal: newMemBufferStatic([]byte(memTotalData1)),
		fused:  newMemBufferStatic([]byte(memUsedData1)),
		used:   &atomic.Float64{},
		dur:    10 * time.Millisecond,
		done:   done,
	}
	memDone.init()
	close(done)

	memOk := &MemCG2Lazy{
		ftotal: newMemBufferStatic([]byte(memTotalData1)),
		fused:  newMemBufferStatic([]byte(memUsedData1)),
		used:   &atomic.Float64{},
		dur:    10 * time.Millisecond,
	}
	memOk.init()

	memFail := &MemCG2Lazy{
		ftotal: newMemBufferStatic([]byte("")),
		fused:  newMemBufferStatic([]byte("")),
		used:   &atomic.Float64{},
		dur:    100 * time.Millisecond,
	}
	memFail.init()

	time.Sleep(20 * time.Millisecond)

	type fields struct {
		mem *MemCG2Lazy
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "ok",
			fields: fields{
				mem: memOk,
			},
			want: 1.4140625,
		},
		{
			name: "done",
			fields: fields{
				mem: memDone,
			},
			want: -2.0,
		},
		{
			name: "fail",
			fields: fields{
				mem: memFail,
			},
			want: -1.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.mem.Used(); got != tt.want {
				t.Errorf("MemCG2Lazy.Used() = %v, want %v", got, tt.want)
			}
		})
	}
}
