package server

import (
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/emicklei/go-restful"
	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/vortex/src/entity"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/server/backend"
)

func globalLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	logger.Infof("%s %s", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
}

func validateTokenMiddleware(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	token, err := request.ParseFromRequest(req.Request, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(backend.SecretKey), nil
		})

	if err == nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// save user ID to requests attributes
			req.SetAttribute("UserID", claims["sub"])
			// save role to requests attributes
			req.SetAttribute("Role", claims["role"])
			chain.ProcessFilter(req, resp)
		} else {
			resp.WriteHeaderAndEntity(http.StatusUnauthorized,
				response.ActionResponse{
					Error:   true,
					Message: "Token is invalid",
				})
			return
		}
	} else {
		logger.Infof("Unauthorized access to this resource")
		resp.WriteHeaderAndEntity(http.StatusUnauthorized,
			response.ActionResponse{
				Error:   true,
				Message: "Unauthorized access to this resource",
			})
		return
	}
}

func rootRole(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	role, ok := req.Attribute("Role").(string)
	if ok && role == entity.RootRole {
		chain.ProcessFilter(req, resp)
	} else {
		log.Printf("User role: %s has no root role: Forbidden", role)
		resp.WriteHeaderAndEntity(http.StatusForbidden,
			response.ActionResponse{
				Error:   true,
				Message: "Permission denied",
			})
		return
	}
}

func userRole(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	role, ok := req.Attribute("Role").(string)
	if ok && role == entity.RootRole || role == entity.UserRole {
		chain.ProcessFilter(req, resp)
	} else {
		log.Printf("User role: %s has no root role: Forbidden", role)
		resp.WriteHeaderAndEntity(http.StatusForbidden,
			response.ActionResponse{
				Error:   true,
				Message: "Permission denied",
			})
		return
	}
}

func guestRole(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	role, ok := req.Attribute("Role").(string)
	if ok && role == entity.RootRole || role == entity.UserRole || role == entity.GuestRole {
		chain.ProcessFilter(req, resp)
	} else {
		log.Printf("User role: %s has no root role: Forbidden", role)
		resp.WriteHeaderAndEntity(http.StatusForbidden,
			response.ActionResponse{
				Error:   true,
				Message: "Permission denied",
			})
		return
	}
}
