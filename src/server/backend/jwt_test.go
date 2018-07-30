package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	roles := "admin"
	tokenString, err := GenerateToken("234243353535330", roles)
	assert.NotNil(t, tokenString)
	assert.NoError(t, err)
}
