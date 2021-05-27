package utils

import "testing"

func Test_CPUPercent(t *testing.T) {
	type args struct {
		lastused  uint64
		used      uint64
		lasttotal uint64
		total     uint64
	}
	tests := []struct {
		name  string
		args  args
		wantP float64
	}{
		{
			name:  "empty",
			args:  args{},
			wantP: 0,
		},
		{
			name: "used zero",
			args: args{
				lastused:  10,
				used:      10,
				lasttotal: 100,
				total:     110,
			},
			wantP: 0,
		},
		{
			name: "total zero",
			args: args{
				lastused:  9,
				used:      10,
				lasttotal: 100,
				total:     100,
			},
			wantP: 100,
		},
		{
			name: "total < 1",
			args: args{
				lastused:  10,
				used:      11,
				lasttotal: 100,
				total:     10,
			},
			wantP: 100,
		},
		{
			name: "used < 1",
			args: args{
				lastused:  12,
				used:      11,
				lasttotal: 100,
				total:     110,
			},
			wantP: 0,
		},
		{
			name: "noramal",
			args: args{
				lastused:  1191919,
				used:      1191920,
				lasttotal: 31900447,
				total:     31900455,
			},
			wantP: 12.5,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if gotP := CPUPercent(tt.args.lastused, tt.args.used, tt.args.lasttotal, tt.args.total); gotP != tt.wantP {
				t.Errorf("percent() = %v, want %v", gotP, tt.wantP)
			}
		})
	}
}
