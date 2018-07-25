package server

import (
	"strings"

	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/web"
)

// TODO move to entity
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginHandler(ctx *web.Context) {
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

	resp.WriteEntity(response.ActionResponse{
		Error:   false,
		Message: tokenString,
	})
}
