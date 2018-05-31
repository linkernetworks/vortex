package server

import (
	"github.com/emicklei/go-restful"
	"github.com/linkernetworks/logger"
	oauth "github.com/linkernetworks/oauth/entity"
)

type ActionResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type SignInResponse struct {
	Error        bool            `json:"error"`
	AuthRequired bool            `json:"authenRequired,omitempty"`
	Message      string          `json:"message"`
	SignInUrl    string          `json:"signInUrl,omitempty"`
	Session      SessionResponse `json:"session,omitempty"`
}

type SessionResponse struct {
	ID          string     `json:"id,omitempty"`
	Token       string     `json:"token,omitempty"`
	ExpiredAt   int64      `json:"expiredAt,omitempty"`
	CurrentUser oauth.User `json:"currentUser,omitempty"`
}

func WriteResponse(r *restful.Response, httpStatus int, res ActionResponse) error {
	return r.WriteHeaderAndEntity(httpStatus, res)
}

func responseErrorWithStatus(r *restful.Response, httpStatus int, msg string) error {
	logger.Error(msg)
	return r.WriteHeaderAndEntity(
		httpStatus,
		ActionResponse{
			Error:   true,
			Message: msg,
		})
}
