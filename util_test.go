package revcatgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	s := []string{"a", "b", "c"}

	assert.True(t, contains(s, "a"))
	assert.False(t, contains(s, "d"))
}
