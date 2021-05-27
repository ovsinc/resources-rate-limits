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
