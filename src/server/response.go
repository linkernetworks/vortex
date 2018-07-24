package server

import (
	"github.com/emicklei/go-restful"
	"github.com/linkernetworks/logger"
	oauth "github.com/linkernetworks/oauth/entity"
)

// ActionResponse is the structure for Response action
type ActionResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

// SignInResponse is the structure for Sign in Response
type SignInResponse struct {
	Error        bool            `json:"error"`
	AuthRequired bool            `json:"authenRequired,omitempty"`
	Message      string          `json:"message"`
	SignInURL    string          `json:"signInURL,omitempty"`
	Session      SessionResponse `json:"session,omitempty"`
}

// SessionResponse is the sreucture for Response session
type SessionResponse struct {
	ID          string     `json:"id,omitempty"`
	Token       string     `json:"token,omitempty"`
	ExpiredAt   int64      `json:"expiredAt,omitempty"`
	CurrentUser oauth.User `json:"currentUser,omitempty"`
}

// WriteResponse will write response
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
