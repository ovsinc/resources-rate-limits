package resourcesratelimits

import (
	"testing"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"

	"github.com/stretchr/testify/assert"
)

func Test_check(t *testing.T) {
	type args struct {
		files []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := check(tt.args.files...); got != tt.want {
				t.Errorf("check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	c := Check()
	assert.Greater(t, c.Type(), rescommon.ResourceType_UNKNOWN)
	assert.Less(t, c.Type(), rescommon.ResourceType_ENDS)
}
