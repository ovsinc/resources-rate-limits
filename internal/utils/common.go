package utils

import (
	"bytes"
	"errors"
	"io"
	"strconv"
)

const (
	CG2Max = "max"

	readBytesLen = 10 // int64 size + '0x0a' + '0x00'
)

var ErrMax = errors.New("parse with keyworld 'max'")

func CPUPercent(lastused, used, lasttotal, total uint64) (p float64) {
	switch {
	case used <= lastused:
		p = 0.0
	case total <= lasttotal:
		p = 100.0
	default:
		p = Percent(float64(used-lastused), float64(total-lasttotal))
	}
	return p
}

func Percent(used float64, total float64) float64 {
	if total == 0 {
		return 0
	}
	return used / total * 100
}

// refer to https://github.com/containerd/cgroups/blob/318312a373405e5e91134d8063d04d59768a1bff/utils.go#L251
func ParseUint(b []byte) (uint64, error) {
	s := string(b)

	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		// check is MAX
		if bytes.HasSuffix(b, []byte(CG2Max)) {
			return 0, ErrMax
		}

		intValue, intErr := strconv.ParseInt(s, 10, 64)
		// 1. Handle negative values greater than MinInt64 (and)
		// 2. Handle negative values lesser than MinInt64
		if intErr == nil && intValue < 0 {
			return 0, nil
		} else if intErr != nil &&
			errors.Is(intErr, strconv.ErrRange) &&
			intValue < 0 {
			return 0, nil
		}

		return 0, err
	}

	return v, nil
}

func ReadUintFromF(f io.Reader) (uint64, error) {
	v := make([]byte, readBytesLen)

	if _, err := f.Read(v); !(err == nil || errors.Is(err, io.EOF)) {
		return 0, err
	}

	return ParseUint(bytes.TrimRight(v, "\n\r \x00"))
}
