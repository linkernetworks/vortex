package server

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/emicklei/go-restful"
	"github.com/linkernetworks/logger"
)

// FIXME using ldconfig go build to give a secretkey
const (
	SecretKey = "linkernetworks"
)

func globalLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	logger.Infof("%s %s", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
}

func validateTokenMiddleware(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	token, err := request.ParseFromRequest(req.Request, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

	if err == nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// save to requests attributes
			req.SetAttribute("UserID", claims["sub"])
			chain.ProcessFilter(req, resp)
		} else {
			resp.WriteHeader(http.StatusUnauthorized)
			logger.Infof("Token is not valid")
		}
	} else {
		resp.WriteHeader(http.StatusUnauthorized)
		logger.Infof("Unauthorized access to this resource")
	}
}
