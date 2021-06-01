package os

import (
	"bytes"
	"io"
)

const (
	data1 = `cpu  695636 4341 199895 27915902 7970 31786 16794 0 0 0
	cpu0 34184 339 10252 1296456 353 1077 1570 0 0 0
	cpu1 37919 380 10769 1294946 378 1185 693 0 0 0
	cpu2 38275 395 10732 1292207 451 1399 2039 0 0 0
	cpu3 38566 372 10888 1294215 423 1013 535 0 0 0
	cpu4 38375 321 10814 1294723 442 1081 430 0 0 0
	cpu5 30476 299 11449 1293194 401 6597 1461 0 0 0
	cpu6 40492 206 11235 1291149 407 1166 647 0 0 0
	cpu7 38739 399 11198 1294028 442 1016 379 0 0`
	data0 = ""

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

type mocStatic struct {
	data []byte
	done bool
}

func (r *mocStatic) Read(p []byte) (n int, err error) {
	if r.done {
		return 0, io.EOF
	}
	buf := bytes.NewBuffer(r.data)
	r.done = true
	return buf.Read(p)
}

func (r *mocStatic) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (r *mocStatic) Close() error { return nil }

func newMocStatic(data []byte) io.ReadSeekCloser {
	return &mocStatic{data: data}
}
