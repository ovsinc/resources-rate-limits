package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_resourceConfig_Type(t *testing.T) {
	type fields struct {
		rtype ResourceType
	}
	tests := []struct {
		name   string
		fields fields
		want   ResourceType
	}{
		{
			name: "nil",
			want: 0,
		},
		{
			name: "os",
			fields: fields{
				rtype: ResourceType_OS,
			},
			want: ResourceType_OS,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			rc := &resourceConfig{
				rtype: tt.fields.rtype,
			}
			if got := rc.Type(); got != tt.want {
				t.Errorf("resourceConfig.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourceConfig_Init(t *testing.T) {
	cnf := NewResourceConfig(0, "/proc/cpuinfo", "/proc/stat")
	assert.NotNil(t, cnf)
	defer cnf.Stop()
	assert.Nil(t, cnf.Init())
	assert.Nil(t, cnf.Init())
	assert.NotNil(t, cnf.File("/proc/cpuinfo"))

	cnf2 := NewResourceConfig(0, "/s/dsdsdsd", "/sds/hyuyu/yuyu")
	assert.NotNil(t, cnf2)
	defer cnf2.Stop()

	assert.NotNil(t, cnf2.Init())
	assert.Nil(t, cnf2.File("/proc/cpuinfo"))
}
