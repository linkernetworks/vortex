package backend

import (
	"testing"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	user := entity.User{
		LoginCredential: entity.LoginCredential{
			Username: "admin@linkernetworks.com",
		},
		DisplayName: "admin",
		Role:        "root",
		FirstName:   "john",
		LastName:    "lin",
		PhoneNumber: "123456789",
	}
	tokenString, err := GenerateToken("234243353535330", user)
	assert.NotNil(t, tokenString)
	assert.NoError(t, err)
}

func TestVerifyToken(t *testing.T) {
	user := entity.User{
		LoginCredential: entity.LoginCredential{
			Username: "admin@linkernetworks.com",
		},
		DisplayName: "admin",
		Role:        "root",
		FirstName:   "john",
		LastName:    "lin",
		PhoneNumber: "123456789",
	}
	tokenString, err := GenerateToken("234243353535330", user)
	assert.NotNil(t, tokenString)
	assert.NoError(t, err)

	isValid := VerifyToken([]byte(tokenString))
	assert.True(t, isValid)
}

func TestInValidVerifyToken(t *testing.T) {
	tokenString := "fakeToken"
	isValid := VerifyToken([]byte(tokenString))
	assert.False(t, isValid)

	tokenString2 := "eyJhbGciOiJIUzI1NiIaInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	isValid2 := VerifyToken([]byte(tokenString2))
	assert.False(t, isValid2)
}
