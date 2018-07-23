package server

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/web"
)

// TODO move to entity
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func authenticateHandler(ctx *web.Context) {
	_, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	user := UserCredentials{}
	if err := req.ReadEntity(&user); err != nil {
		response.Forbidden(req.Request, resp.ResponseWriter, err)
		return
	}

	// TODO query mongodb to check user account and password
	if strings.ToLower(user.Username) != "someone" || user.Password != "p@ssword" {
		response.Forbidden(req.Request, resp.ResponseWriter)
		return
	}

	tokenString, err := generateToken(strings.ToLower(user.Username))
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	resp.WriteEntity(ActionResponse{
		Error:   false,
		Message: tokenString,
	})
}

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
