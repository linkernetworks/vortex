package http

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// ErrorPayload is the Structure to contain the Error Message from the HTTP response
// Contains the Error Messagess (At most two errors)
type ErrorPayload struct {
	XMLName         xml.Name `json:"-" xml:"response"`
	Error           bool     `json:"error" xml:"error"`
	Message         string   `json:"message" xml:"message"`
	PreviousMessage string   `json:"previousMessage,omitempty" xml:"previousMessage,omitempty"`
}

// ActionResponse is the structure for Response action
type ActionResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

// NewErrorPayload will return the ErrorPayload message according the parameters (errors)
// The ErrorPayload contains at most error messages
func NewErrorPayload(errs ...error) ErrorPayload {
	payload := ErrorPayload{Error: true}

	if len(errs) > 0 {
		payload.Message = errs[0].Error()
	}
	if len(errs) > 1 {
		payload.PreviousMessage = errs[1].Error()
	}

	return payload
}

// EncodeErrorPayload is the function to get the error from the request and ErrorPayload
// It should contains the errorType and error content
func EncodeErrorPayload(req *http.Request, payload ErrorPayload) (out []byte, contentType string, encErr error) {
	contentType = req.Header.Get("content-type")

	switch contentType {
	case "application/json", "text/json":
		out, encErr = json.Marshal(payload)
	case "application/xml", "text/xml":
		out, encErr = xml.Marshal(payload)
	default:
		contentType = "application/json"
		out, encErr = json.Marshal(payload)
	}

	return out, contentType, encErr
}

// WriteStatusAndError will write the error status and error code to the http request and return the written byte count of the http request
func WriteStatusAndError(req *http.Request, resp http.ResponseWriter, status int, errs ...error) (int, error) {
	payload := NewErrorPayload(errs...)
	out, contentType, encErr := EncodeErrorPayload(req, payload)
	if encErr != nil {
		resp.Write([]byte("failed to encode payload"))
		return 0, encErr
	}

	resp.WriteHeader(status)
	resp.Header().Set("Content-Type", contentType)
	return resp.Write(out)
}

// Forbidden will set the status code http.StatusForbidden to the HTTP response message
func Forbidden(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusForbidden, errs...)
}

// BadRequest will set the status code http.StatusBadRequest to the HTTP response message
func BadRequest(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusBadRequest, errs...)
}

// OK will set the status code http.StatusOK to the HTTP response message
func OK(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusOK, errs...)
}

// NotFound will set the status code http.StatusNotFound to the HTTP response message
func NotFound(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusNotFound, errs...)
}

// Unauthorized will set the status code http.StatusUnauthorized to the HTTP response message
func Unauthorized(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusUnauthorized, errs...)
}

// InternalServerError will set the status code http.StatusInternalServerError to the HTTP response message
func InternalServerError(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusInternalServerError, errs...)
}

// Conflict will set the status code http.StatusConflict to the HTTP response message
func Conflict(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusConflict, errs...)
}

// MethodNotAllow will set the status code http.StatusMethodNotAllowed, to the HTTP response message
func MethodNotAllow(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusMethodNotAllowed, errs...)
}

// UnprocessableEntity will set the status code http.StatusUnprocessableEntity to the HTTP response message
func UnprocessableEntity(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusUnprocessableEntity, errs...)
}

// NotAcceptable will set the status code http.StatusNotAcceptable to the HTTP response message
func NotAcceptable(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusNotAcceptable, errs...)
}
