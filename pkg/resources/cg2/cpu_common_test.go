package cg2

import (
	"bytes"
	"io"
	"testing"
)

const (
	cpuTotalEmpty  = ``
	cpuTotalBadMax = `dd 100
`
	cpuTotalBadPer = `100 max
`
	cpuTotalShort = `10
`

	cpuTotal = `50000 100000
`
	cpuTotalMax = `max 100000
`

	cpuStat = `usage_usec 31585
user_usec 21082
system_usec 10503
nr_periods 0
nr_throttled 0
throttled_usec 0
`
	cpuStatBad = `ddd 999
`
	cpuStatParseErr = `usage_usec fff
user_usec 21082
system_usec 10503
nr_periods 0
`
)

type cpuMocStatic struct {
	data []byte
}

func (r *cpuMocStatic) Read(p []byte) (n int, err error) {
	buf := bytes.NewBuffer(r.data)
	return buf.Read(p)
}

func (r *cpuMocStatic) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (r *cpuMocStatic) Close() error { return nil }

func newCPUBufferStatic(data []byte) io.ReadSeekCloser {
	return &cpuMocStatic{data: data}
}

//

const (
	memUsedZero = `
`
	memTotalData1 = `104857600
`
	memTotalDataMax = `max
`
	memUsedData1 = `1482752
`
	meminfo = `MemTotal:       32338608 kB
MemFree:        21588516 kB
MemAvailable:   25582292 kB
Buffers:          320640 kB
Cached:          5661316 kB
SwapCached:            0 kB
Active:          2476796 kB
Inactive:        6881216 kB
Active(anon):     616280 kB
Inactive(anon):  4544412 kB
Active(file):    1860516 kB
Inactive(file):  2336804 kB
Unevictable:      748392 kB
Mlocked:              96 kB
SwapTotal:       6291452 kB
SwapFree:        6291452 kB
Dirty:               600 kB
Writeback:             0 kB
AnonPages:       4120840 kB
Mapped:          1081472 kB
Shmem:           1788240 kB
KReclaimable:     264592 kB
Slab:             418556 kB
SReclaimable:     264592 kB
SUnreclaim:       153964 kB
KernelStack:       26224 kB
PageTables:        79148 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:    22460756 kB
Committed_AS:   28656504 kB
VmallocTotal:   34359738367 kB
VmallocUsed:       63400 kB
VmallocChunk:          0 kB
Percpu:            11904 kB
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
DirectMap4k:      465144 kB
DirectMap2M:    14721024 kB
DirectMap1G:    17825792 kB
`
)

type memMocStatic struct {
	data []byte
}

func (r *memMocStatic) Read(p []byte) (n int, err error) {
	buf := bytes.NewBuffer(r.data)
	return buf.Read(p)
}

func (r *memMocStatic) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (r *memMocStatic) Close() error { return nil }

func newMemBufferStatic(data []byte) io.ReadSeekCloser {
	return &memMocStatic{data: data}
}

func Test_parceCG2MaxCPUUintFromF(t *testing.T) {
	type args struct {
		f io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantPeriod uint64
		wantMax    uint64
		wantErr    bool
	}{
		{
			name: "empty",
			args: args{
				f: newCPUBufferStatic([]byte(cpuTotalEmpty)),
			},
			wantErr: true,
		},
		{
			name: "short",
			args: args{
				f: newCPUBufferStatic([]byte(cpuTotalShort)),
			},
			wantErr: true,
		},
		{
			name: "bad max",
			args: args{
				f: newCPUBufferStatic([]byte(cpuTotalBadMax)),
			},
			wantErr: true,
		},
		{
			name: "bad preiod",
			args: args{
				f: newCPUBufferStatic([]byte(cpuTotalBadPer)),
			},
			wantErr: true,
		},
		{
			name: "max",
			args: args{
				f: newCPUBufferStatic([]byte(cpuTotalMax)),
			},
			wantPeriod: 100000,
			wantMax:    100000,
		},
		{
			name: "ok",
			args: args{
				f: newCPUBufferStatic([]byte(cpuTotal)),
			},
			wantPeriod: 100000,
			wantMax:    50000,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			gotPeriod, gotMax, err := parceCG2MaxCPUUintFromF(tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParceCG2CPUUintFromF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPeriod != tt.wantPeriod {
				t.Errorf("ParceCG2CPUUintFromF() gotPeriod = %v, want %v", gotPeriod, tt.wantPeriod)
			}
			if gotMax != tt.wantMax {
				t.Errorf("ParceCG2CPUUintFromF() gotMax = %v, want %v", gotMax, tt.wantMax)
			}
		})
	}
}

func Test_parceCG2StatCPUUintFromF(t *testing.T) {
	type args struct {
		f io.Reader
	}
	tests := []struct {
		name     string
		args     args
		wantUsed uint64
		wantErr  bool
	}{
		{
			name: "empty",
			args: args{
				f: newCPUBufferStatic([]byte("")),
			},
			wantErr: true,
		},
		{
			name: "bad",
			args: args{
				f: newCPUBufferStatic([]byte(cpuStatBad)),
			},
			wantErr: true,
		},
		{
			name: "parse err",
			args: args{
				f: newCPUBufferStatic([]byte(cpuStatParseErr)),
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				f: newCPUBufferStatic([]byte(cpuStat)),
			},
			wantUsed: 31585,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			gotUsed, err := parceCG2StatCPUUintFromF(tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("parceCG2StatCPUUintFromF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUsed != tt.wantUsed {
				t.Errorf("parceCG2StatCPUUintFromF() = %v, want %v", gotUsed, tt.wantUsed)
			}
		})
	}
}
