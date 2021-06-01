package cg2

import (
	"io"

	"github.com/ovsinc/errors"

	"github.com/ovsinc/resources-rate-limits/internal/utils"
	"github.com/ovsinc/resources-rate-limits/pkg/resources/os"
)

/*
$ cat /sys/fs/cgroup/memory.max
max

$ cat /sys/fs/cgroup/memory.max
104857600

$ cat /sys/fs/cgroup/memory.current
528384
*/

func getMemInfo(ftotal, fused, procmem io.ReadSeeker) (uint64, uint64, error) {
	_, err := ftotal.Seek(0, 0)
	if err != nil {
		return 0, 0, err
	}

	_, err = fused.Seek(0, 0)
	if err != nil {
		return 0, 0, err
	}

	var total, used uint64

	total, err = utils.ReadUintFromF(ftotal)
	switch {
	case errors.Is(err, utils.ErrMax):
		total, _, err = os.GetMemInfo(procmem) // in kb
		if err != nil {
			return 0, 0, err
		}
		total = total * 1024
	case err != nil:
		return 0, 0, err
	}

	used, err = utils.ReadUintFromF(fused)
	if err != nil {
		return 0, 0, err
	}

	return total, used, nil
}
