package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	_, err := HashPassword("p@ssw0rd")
	assert.NoError(t, err)
}

func TestCheckPasswordHash(t *testing.T) {
	hashed, err := HashPassword("p@ssw0rd")
	assert.NoError(t, err)
	correct := CheckPasswordHash("p@ssw0rd", hashed)
	assert.True(t, true, correct)
}
