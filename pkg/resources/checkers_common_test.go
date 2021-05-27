package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	c := Check()
	assert.Greater(t, c, ResourceType_UNKNOWN)
	assert.Less(t, c, ResourceType_ENDS)
}
