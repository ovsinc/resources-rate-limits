package cg1

import (
	"bytes"
	"io"
)

const (
	CPUtotal = `100000
`
	CPUused = `53704441
`

	CPUtotalUnquoted = `-1
`

	CPUtotalFail = ``
	CPUusedFail  = ``

	MemTotal = `10485760
`
	MemUsed = `294912
`

	MemTotalUnquoted = `-1
`
	MemTotalFail = ``
	MemUsedFail  = ``
)

type mocStatic struct {
	data []byte
}

func (r *mocStatic) Read(p []byte) (n int, err error) {
	buf := bytes.NewBuffer(r.data)
	return buf.Read(p)
}

func (r *mocStatic) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (r *mocStatic) Close() error { return nil }

func newMocStatic(data []byte) io.ReadSeekCloser {
	return &mocStatic{data: data}
}
