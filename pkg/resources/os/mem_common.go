package os

import (
	"bufio"
	"bytes"
	"errors"
	"io"

	"github.com/ovsinc/resources-rate-limits/internal/utils"
)

/*
$ cat /proc/meminfo
MemTotal:       32336560 kB
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
DirectMap1G:    20971520 kB
*/

var (
	ErrUnknownScanErr = errors.New("scan failed")
	ErrNoMemInfoFile  = errors.New("not /proc/meminfo file")
)

func getMemInfo(f io.ReadSeeker) (uint64, uint64, error) {
	_, err := f.Seek(0, 0)
	if err != nil {
		return 0, 0, err
	}

	var (
		total, used uint64

		free, buffers, cached uint64
		sReclaimable          uint64
	)

	scanner := bufio.NewScanner(f)

SCANLOOP:
	for scanner.Scan() {
		fields := bytes.Split(scanner.Bytes(), []byte(":"))
		if len(fields) != 2 {
			return 0, 0, ErrNoMemInfoFile
		}

		key := bytes.TrimSpace(fields[0])
		value := bytes.ReplaceAll(
			bytes.TrimSpace(fields[1]),
			[]byte(" kB"), []byte(""),
		)

		switch {
		case bytes.Equal([]byte("MemTotal"), key):
			t, err := utils.ParseUint(value)
			if err != nil {
				return 0, 0, err
			}
			total = t

		case bytes.Equal([]byte("MemFree"), key):
			t, err := utils.ParseUint(value)
			if err != nil {
				return 0, 0, err
			}
			free = t

		case bytes.Equal([]byte("Buffers"), key):
			t, err := utils.ParseUint(value)
			if err != nil {
				return 0, 0, err
			}
			buffers = t

		case bytes.Equal([]byte("Cached"), key):
			t, err := utils.ParseUint(value)
			if err != nil {
				return 0, 0, err
			}
			cached = t

		case bytes.Equal([]byte("SReclaimable"), key):
			t, err := utils.ParseUint(value)
			if err != nil {
				return 0, 0, err
			}
			sReclaimable = t

			if free > 0 && buffers > 0 && cached > 0 && sReclaimable > 0 {
				break SCANLOOP
			}
		}
	}

	if total == 0 {
		return 0, 0, ErrUnknownScanErr
	}

	used = total -
		free - cached - sReclaimable - buffers

	return total, used, nil
}
