// +build os

package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCG1Check(t *testing.T) {
	c := Check()
	assert.Equal(t, c, ResourceType_OS)
}
