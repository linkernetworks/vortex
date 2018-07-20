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

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	resp.WriteEntity(ActionResponse{
		Error:   false,
		Message: tokenString,
	})
}
