package cg1

import (
	"testing"

	"github.com/ovsinc/resources-rate-limits/pkg/resources/common/moc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCPUCG1Simple_info_moc(t *testing.T) {
	// ok

	oktotal, err := moc.CreateTemp("/tmp", []byte(CPUtotal))
	require.Nil(t, err)
	defer oktotal.Remove()

	okused, err := moc.CreateTemp("/tmp", []byte(CPUused))
	require.Nil(t, err)
	defer okused.Remove()

	// empty

	emptytotal, err := moc.CreateTemp("/tmp", []byte(""))
	require.Nil(t, err)
	defer emptytotal.Remove()

	emptyused, err := moc.CreateTemp("/tmp", []byte(""))
	require.Nil(t, err)
	defer emptyused.Remove()

	// unqouted

	totunquouted, err := moc.CreateTemp("/tmp", []byte(CPUtotalUnquoted))
	require.Nil(t, err)
	defer totunquouted.Remove()

	type fields struct {
		limit string
		used  string
	}
	tests := []struct {
		name      string
		fields    fields
		wantTotal uint64
		wantUsed  uint64
		wantErr   bool
	}{
		{
			name: "ok",
			fields: fields{
				limit: oktotal.File().Name(),
				used:  okused.File().Name(),
			},
			wantTotal: 100000,
			wantUsed:  53704441,
		},

		{
			name: "empty",
			fields: fields{
				limit: emptytotal.File().Name(),
				used:  emptyused.File().Name(),
			},
			wantErr: true,
		},

		{
			name: "empty total",
			fields: fields{
				limit: emptytotal.File().Name(),
				used:  okused.File().Name(),
			},
			wantErr: true,
		},

		{
			name: "empty used",
			fields: fields{
				limit: oktotal.File().Name(),
				used:  emptyused.File().Name(),
			},
			wantErr: true,
		},

		{
			name: "unquoted",
			fields: fields{
				limit: totunquouted.File().Name(),
				used:  okused.File().Name(),
			},
			wantUsed:  53704441,
			wantTotal: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := &CPUCG1Simple{
				limit: tt.fields.limit,
				used:  tt.fields.used,
			}
			gotTotal, gotUsed, err := cpu.info()

			if (err != nil) != tt.wantErr {
				t.Errorf("CPUCG1Simple.info() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("CPUCG1Simple.info() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
			}
			if gotUsed != tt.wantUsed {
				t.Errorf("CPUCG1Simple.info() gotUsed = %v, want %v", gotUsed, tt.wantUsed)
			}
		})
	}
}

func TestCPUCG1Simple_Used_moc(t *testing.T) {
	// ok

	oktotal, err := moc.CreateTemp("/tmp", []byte(CPUtotal))
	require.Nil(t, err)
	defer oktotal.Remove()

	okused, err := moc.CreateTemp("/tmp", []byte(CPUused))
	require.Nil(t, err)
	defer okused.Remove()

	// empty

	emptytotal, err := moc.CreateTemp("/tmp", []byte(""))
	require.Nil(t, err)
	defer emptytotal.Remove()

	emptyused, err := moc.CreateTemp("/tmp", []byte(""))
	require.Nil(t, err)
	defer emptyused.Remove()

	// unqouted

	totunquouted, err := moc.CreateTemp("/tmp", []byte(CPUtotalUnquoted))
	require.Nil(t, err)
	defer totunquouted.Remove()

	type fields struct {
		limit string
		used  string
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "ok",
			fields: fields{
				limit: oktotal.File().Name(),
				used:  okused.File().Name(),
			},
			want: 0.053704441000000006,
		},

		{
			name: "empty",
			fields: fields{
				limit: emptytotal.File().Name(),
				used:  emptyused.File().Name(),
			},
			want: -1.0,
		},

		{
			name: "empty total",
			fields: fields{
				limit: emptytotal.File().Name(),
				used:  okused.File().Name(),
			},
			want: -1.0,
		},

		{
			name: "empty used",
			fields: fields{
				limit: oktotal.File().Name(),
				used:  emptyused.File().Name(),
			},
			want: -1.0,
		},

		{
			name: "unquoted",
			fields: fields{
				limit: totunquouted.File().Name(),
				used:  okused.File().Name(),
			},
			want: 0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := &CPUCG1Simple{
				limit: tt.fields.limit,
				used:  tt.fields.used,
			}
			if got := cg.Used(); got != tt.want {
				t.Errorf("CPUCG1Simple.Used() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCPUSimple(t *testing.T) {
	_, err := NewCPUSimple()
	assert.Nil(t, err)
}
