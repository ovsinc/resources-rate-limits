package os

import (
	origos "os"
	"testing"
	"time"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkNewMemLazy_info_mock(b *testing.B) {
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

func TestMemOSLazy_Used_Sys(t *testing.T) {
	f, err := origos.Open(rescommon.RAMFilenameInfoProc)
	require.Nil(t, err)
	defer f.Close()

	done := make(chan struct{})
	defer close(done)

	cnf := rescommon.NewResourceConfig(rescommon.ResourceType_OS, rescommon.RAMFilenameInfoProc)
	require.Nil(t, cnf.Init())
	defer cnf.Stop()

	mi, err := NewMemLazy(done, cnf, 100*time.Millisecond)
	assert.Nil(t, err)

	time.Sleep(250 * time.Millisecond)

	u := mi.Used()
	assert.Greater(t, u, float64(0))
	assert.Less(t, u, float64(100))
}

func TestMemOSLazy_info_mock(t *testing.T) {
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
