package backend

import (
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/vortex/src/entity"
)

// GenerateToken is for generating token
func GenerateToken(userID string, user entity.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims{
		// issuer of the claim
		"exp": time.Now().Add(time.Hour * time.Duration(8)).Unix(),
		// issued-at time
		"iat": time.Now().Unix(),
		// user role
		"role": user.Role,
		// user email
		"username": user.LoginCredential.Username,
		// user display name
		"displayName": user.DisplayName,
		// the subject of this token. This is the user associated with the relevant action
		"sub": userID,
	}
	return token.SignedString([]byte(SecretKey))
}

// VerifyToken is for verifing the JWT
func VerifyToken(tokenData []byte) bool {
	// trim possible whitespace from token
	tokenData = regexp.MustCompile(`\s*$`).ReplaceAll(tokenData, []byte{})

	// Parse the token
	token, err := jwt.Parse(string(tokenData), func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	// Print an error if we can't parse for some reason
	if err != nil {
		logger.Infof("Couldn't parse token: %v", err)
		return false
	}

	// Is token invalid?
	if !token.Valid {
		logger.Infof("Token is invalid")
		return false
	}
	return true
}
