package server

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func generateToken(userUUID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims{
		// issuer of the claim
		"exp": time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		// issued-at time
		"iat": time.Now().Unix(),
		// the subject of this token. This is the user associated with the relevant action
		"sub": userUUID,
	}
	return token.SignedString([]byte(SecretKey))
}
