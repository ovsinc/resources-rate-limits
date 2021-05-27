package cg2

import (
	"bufio"
	"bytes"
	"errors"
	"io"

	"github.com/ovsinc/resources-rate-limits/internal/utils"
)

/*
$ cat /sys/fs/cgroup/cpu.stat
usage_usec 46578
user_usec 31044
system_usec 15534
nr_periods 25
nr_throttled 0
throttled_usec 0

$ cat /sys/fs/cgroup/cpu.max
100000 100000
*/

const (
	cg2usage = "usage_usec"

	readBytesLen = 19 // (2 * int64 size) + '0x0a' + ' ' +  '0x00'
)

var (
	ErrParseCPUUsage       = errors.New("parse cgroups2 cpu usage fails")
	ErrCPUUsageNotUsage    = errors.New("cgroups2 cpu usage not contails `usage_usec`")
	ErrCG2MaxCPUParseFails = errors.New("cgroups2 cpu max fails")
)

func parceCG2StatCPUUintFromF(f io.Reader) (used uint64, err error) {
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	fields := bytes.Fields(bytes.TrimRight(scanner.Bytes(), " "))
	switch {
	case len(fields) != 2:
		return 0, ErrParseCPUUsage
	case !bytes.Equal(fields[0], []byte(cg2usage)):
		return 0, ErrCPUUsageNotUsage
	}

	used, err = utils.ParseUint(fields[1])
	if err != nil {
		return 0, err
	}

	return used, nil
}

func parceCG2MaxCPUUintFromF(f io.Reader) (period, max uint64, err error) {
	v := make([]byte, readBytesLen)

	if _, err := f.Read(v); !(err == nil || errors.Is(err, io.EOF)) {
		return 0, 0, err
	}

	maxper := bytes.Split(bytes.TrimRight(v, "\n\r \x00"), []byte(" "))

	if len(maxper) != 2 {
		return 0, 0, ErrCG2MaxCPUParseFails
	}

	period, err = utils.ParseUint(maxper[1])
	if err != nil {
		return 0, 0, err
	}

	max, err = utils.ParseUint(maxper[0])
	switch {
	case errors.Is(err, utils.ErrMax):
		max = period
	case err != nil:
		return 0, 0, err
	}

	return period, max, nil
}

func getCPUInfo(ftotal, fused io.ReadSeeker) (total uint64, used uint64, err error) {
	_, err = ftotal.Seek(0, 0)
	if err != nil {
		return 0, 0, err
	}

	_, total, err = parceCG2MaxCPUUintFromF(ftotal)
	if err != nil {
		return 0, 0, err
	}

	_, err = fused.Seek(0, 0)
	if err != nil {
		return 0, 0, err
	}

	used, err = parceCG2StatCPUUintFromF(fused)
	if err != nil {
		return 0, 0, err
	}

	return total, used, nil
}
