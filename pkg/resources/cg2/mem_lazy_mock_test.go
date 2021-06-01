package cg2

import (
	"io"
	"testing"
	"time"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
	resmoc "github.com/ovsinc/resources-rate-limits/pkg/resources/common/moc"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func TestNewMemLazy_mock(t *testing.T) {
	done := make(chan struct{})
	defer close(done)

	ftotal := newMemBufferStatic([]byte(memTotalData1))
	fused := newMemBufferStatic([]byte(memUsedData1))
	fprocmem := newMemBufferStatic([]byte(meminfo))

	cnf := &resmoc.ResourceConfigMoc{
		Rtype: rescommon.ResourceType_CG2,
		FF: map[string]io.ReadSeekCloser{
			rescommon.CGroup2MemLimitPath: ftotal,
			rescommon.CGroup2MemUsagePath: fused,
			rescommon.RAMFilenameInfoProc: fprocmem,
		},
	}

	_, err := NewMemLazy(done, cnf, time.Millisecond)
	assert.Nil(t, err)

	_, err = NewMemLazy(done, cnf, 0)
	assert.NotNil(t, err)

	_, err = NewMemLazy(done, nil, time.Millisecond)
	assert.NotNil(t, err)
}

func TestMemCG2Lazy_info_mock(t *testing.T) {
	type fields struct {
		ftotal   io.ReadSeekCloser
		fused    io.ReadSeekCloser
		fprocmem io.ReadSeekCloser
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
			wantTotal: 32338608,
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

func TestMemCG2Lazy_Used_mock(t *testing.T) {
	done := make(chan struct{})
	defer close(done)

	mem := &MemCG2Lazy{
		ftotal: newMemBufferStatic([]byte(memTotalData1)),
		fused:  newMemBufferStatic([]byte(memUsedData1)),
		used:   &atomic.Float64{},
		dur:    100 * time.Millisecond,
	}

	mem.init()

	time.Sleep(time.Second)

	u := mem.Used()
	assert.Greater(t, u, float64(0))
	assert.Less(t, u, float64(100))
}
