package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateVethName(t *testing.T) {
	o := SHA256String("12345678")
	assert.Equal(t, "ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f", o)
}
