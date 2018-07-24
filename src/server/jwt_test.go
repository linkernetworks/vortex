package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	tokenString, err := generateToken("234243353535330")
	assert.NotNil(t, tokenString)
	assert.NoError(t, err)
}
