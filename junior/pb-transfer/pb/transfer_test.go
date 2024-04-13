package pb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransfer(t *testing.T) {
	assert.Equal(t, Transfer(1), 1000)
	assert.Equal(t, Transfer(10), 100)
	assert.Equal(t, Transfer(100), 10)
	assert.Equal(t, Transfer(500), 2)
	assert.Equal(t, Transfer(1000), 1)
}
