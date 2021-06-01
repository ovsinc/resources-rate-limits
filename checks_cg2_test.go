// +build cg2

package resourcesratelimits

import (
	"testing"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"

	"github.com/stretchr/testify/assert"
)

func TestCG1Check(t *testing.T) {
	c := Check()
	assert.Equal(t, c.Type(), rescommon.ResourceType_CG2)
}
