package backend

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateToken is for generating token
func GenerateToken(userUUID string, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims{
		// issuer of the claim
		"exp": time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		// issued-at time
		"iat": time.Now().Unix(),
		// user role
		"role": role,
		// the subject of this token. This is the user associated with the relevant action
		"sub": userUUID,
	}
	return token.SignedString([]byte(SecretKey))
}
