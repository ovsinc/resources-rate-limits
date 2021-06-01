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

func BenchmarkNewMemLazy_info_moc(b *testing.B) {
	mi := &MemOSLazy{
		f: newMocStatic([]byte(memdata1)),
	}

	_, _, err := mi.info()
	require.Nil(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = mi.info()
	}
}

//

func TestMemOSLazy_info_moc(t *testing.T) {
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
			name: "moc static reader",
			fields: fields{
				f: newMocStatic([]byte(memdata1)),
			},
			wantTotal: 32336560,
			wantUsed:  3641908,
		},
		{
			name: "empty",
			fields: fields{
				f: newMocStatic([]byte("")),
			},
			wantErr: true,
		},
		{
			name: "bad file",
			fields: fields{
				f: newMocStatic([]byte("MemTotal:       kkk kB")),
			},
			wantErr: true,
		},
		{
			name: "bad file 2",
			fields: fields{
				f: newMocStatic([]byte("fff       kkk kB")),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mem := &MemOSLazy{
				f: tt.fields.f,
			}
			total, used, err := mem.info()
			if (err != nil) != tt.wantErr {
				t.Errorf("OSmem.getMemInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if total != tt.wantTotal || used != tt.wantUsed {
				t.Errorf("OSmem.getMemInfo() = %d | %d, want %d | %d", used, total, tt.wantUsed, tt.wantTotal)
			}
		})
	}
}

func TestMemOSLazy_Used_moc(t *testing.T) {
	memOk := &MemOSLazy{
		f:    newMocStatic([]byte(memdata1)),
		dur:  10 * time.Millisecond,
		used: &atomic.Float64{},
	}
	memOk.init()

	done := make(chan struct{})
	memDone := &MemOSLazy{
		f:    newMocStatic([]byte(memdata1)),
		dur:  100 * time.Millisecond,
		used: &atomic.Float64{},
		done: done,
	}
	memDone.init()
	close(done)

	memFail := &MemOSLazy{
		f:    newMocStatic([]byte("")),
		dur:  10 * time.Millisecond,
		used: &atomic.Float64{},
	}
	memFail.init()

	time.Sleep(20 * time.Millisecond)

	type fields struct {
		mem *MemOSLazy
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
			want: 11.262509060951444,
		},
		{
			name: "fail",
			fields: fields{
				mem: memFail,
			},
			want: -1.0,
		},
		{
			name: "done",
			fields: fields{
				mem: memDone,
			},
			want: -2.0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.mem.Used(); got != tt.want {
				t.Errorf("MemOSLazy.Used() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMemLazy(t *testing.T) {
	done := make(chan struct{})
	defer close(done)

	cnf := resmoc.ResourceConfigMoc{
		Rtype: rescommon.ResourceType_OS,
		FF: map[string]rescommon.ReadSeekCloser{
			rescommon.RAMFilenameInfoProc: newMocStatic([]byte(memdata1)),
		},
	}

	mem, err := NewMemLazy(done, &cnf, 100*time.Millisecond)
	require.Nil(t, err)
	assert.Equal(t, mem.Used(), 11.262509060951444)

	_, err = NewMemLazy(done, nil, 100*time.Millisecond)
	assert.Error(t, err, rescommon.ErrNoResourceConfig)

	_, err = NewMemLazy(done, &cnf, 0)
	assert.Error(t, err, rescommon.ErrTickPeriodZero)

	cnf = resmoc.ResourceConfigMoc{
		Rtype: rescommon.ResourceType_OS,
		FF:    map[string]rescommon.ReadSeekCloser{},
	}

	_, err = NewMemLazy(done, &cnf, 100*time.Millisecond)
	assert.NotNil(t, err)
}
