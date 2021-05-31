package fiber

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_defaultConfig(t *testing.T) {
	cnf := defaultConfig(Config{})
	assert.Equal(t, cnf, DefaultConfig)
}
