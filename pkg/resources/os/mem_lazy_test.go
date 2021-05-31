package os

import (
	"bytes"
	"io"
	origos "os"
	"testing"
	"time"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	memdata1 = `MemTotal:       32336560 kB
MemFree:        23431752 kB
MemAvailable:   27334948 kB
Buffers:          627952 kB
Cached:          4254340 kB
SwapCached:            0 kB
Active:          1889392 kB
Inactive:        5636988 kB
Active(anon):      43824 kB
Inactive(anon):  3491828 kB
Active(file):    1845568 kB
Inactive(file):  2145160 kB
Unevictable:      632008 kB
Mlocked:              64 kB
SwapTotal:       6291452 kB
SwapFree:        6291452 kB
Dirty:               560 kB
Writeback:             0 kB
AnonPages:       3274160 kB
Mapped:          1129628 kB
Shmem:            895168 kB
KReclaimable:     380608 kB
Slab:             531732 kB
SReclaimable:     380608 kB
SUnreclaim:       151124 kB
KernelStack:       24416 kB
PageTables:        73980 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:    22459732 kB
Committed_AS:   25683884 kB
VmallocTotal:   34359738367 kB
VmallocUsed:       62432 kB
VmallocChunk:          0 kB
Percpu:            10496 kB
HardwareCorrupted:     0 kB
AnonHugePages:         0 kB
ShmemHugePages:        0 kB
ShmemPmdMapped:        0 kB
FileHugePages:         0 kB
FilePmdMapped:         0 kB
CmaTotal:              0 kB
CmaFree:               0 kB
HugePages_Total:       0
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
Hugetlb:               0 kB
DirectMap4k:      444664 kB
DirectMap2M:    11595776 kB
DirectMap1G:    20971520 kB`
)

type memMocStatic struct {
	data []byte
	done bool
}

func (r *memMocStatic) Read(p []byte) (n int, err error) {
	if r.done {
		return 0, io.EOF
	}
	buf := bytes.NewBuffer(r.data)
	r.done = true
	return buf.Read(p)
}

func (r *memMocStatic) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (r *memMocStatic) Close() error { return nil }

func newMemBufferStatic(data []byte) io.ReadSeekCloser {
	return &memMocStatic{
		data: data,
	}
}

func BenchmarkNewMemLazy_info_mock(b *testing.B) {
	mi := &MemOSLazy{
		f: newMemBufferStatic([]byte(memdata1)),
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
			name: "moc static reader",
			fields: fields{
				f: newMemBufferStatic([]byte(memdata1)),
			},
			wantTotal: 32336560,
			wantUsed:  3641908,
		},
		{
			name: "empty",
			fields: fields{
				f: newMemBufferStatic([]byte("")),
			},
			wantErr: true,
		},
		{
			name: "bad file",
			fields: fields{
				f: newMemBufferStatic([]byte("MemTotal:       kkk kB")),
			},
			wantErr: true,
		},
		{
			name: "bad file 2",
			fields: fields{
				f: newMemBufferStatic([]byte("fff       kkk kB")),
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
