package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	_, err := hashPassword("p@ssw0rd")
	assert.NoError(t, err)
}

func TestCheckPasswordHash(t *testing.T) {
	hashed, err := hashPassword("p@ssw0rd")
	assert.NoError(t, err)
	correct := checkPasswordHash("p@ssw0rd", hashed)
	assert.True(t, true, correct)
}
