package cg1

import (
	"io"

	"github.com/ovsinc/resources-rate-limits/internal/utils"
)

func readInfo(f io.ReadSeeker) (uint64, error) {
	_, err := f.Seek(0, 0)
	if err != nil {
		return 0, err
	}
	return utils.ReadUintFromF(f)
}
