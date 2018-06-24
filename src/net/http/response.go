package http

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// The Structure to contain the Error Message from the HTTP response
// Contains the Error Messagess (At most two errors)
type ErrorPayload struct {
	XMLName         xml.Name `json:"-" xml:"response"`
	Error           bool     `json:"error" xml:"error"`
	Message         string   `json:"message" xml:"message"`
	PreviousMessage string   `json:"previousMessage,omitempty" xml:"previousMessage,omitempty"`
}

// Return the ErrorPayload message according the parameters (errors)
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

// The function to get the error from the request and ErrorPayload
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

// Write the error status and error code to the http request and return the written byte count of the http request
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

// Set the status code http.StatusForbidden to the HTTP response message
func Forbidden(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusForbidden, errs...)
}

// Set the status code http.StatusBadRequest to the HTTP response message
func BadRequest(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusBadRequest, errs...)
}

// Set the status code http.StatusOK to the HTTP response message
func OK(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusOK, errs...)
}

// Set the status code http.StatusNotFound to the HTTP response message
func NotFound(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusNotFound, errs...)
}

// Set the status code http.StatusUnauthorized to the HTTP response message
func Unauthorized(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusUnauthorized, errs...)
}

// Set the status code http.StatusInternalServerError to the HTTP response message
func InternalServerError(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusInternalServerError, errs...)
}

// Set the status code http.StatusConflict to the HTTP response message
func Conflict(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusConflict, errs...)
}

// Set the status code http.StatusUnprocessableEntity to the HTTP response message
func UnprocessableEntity(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusUnprocessableEntity, errs...)
}
