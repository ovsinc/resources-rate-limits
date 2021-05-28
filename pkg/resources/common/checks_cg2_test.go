// +build cg2

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCG1Check(t *testing.T) {
	c := Check()
	assert.Equal(t, c, ResourceType_CG2)
}
