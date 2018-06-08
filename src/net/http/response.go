package http

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type ErrorPayload struct {
	XMLName         xml.Name `json:"-" xml:"response"`
	Error           bool     `json:"error" xml:"error"`
	Message         string   `json:"message" xml:"message"`
	PreviousMessage string   `json:"previousMessage,omitempty" xml:"previousMessage,omitempty"`
}

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

func Forbidden(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusForbidden, errs...)
}

func BadRequest(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusBadRequest, errs...)
}

func OK(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusOK, errs...)
}

// Response http.StatusNotFound
func NotFound(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusNotFound, errs...)
}

func Unauthorized(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusUnauthorized, errs...)
}

// Response http.StatusInternalServerError
func InternalServerError(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusInternalServerError, errs...)
}

// Response http.StatusConflict
func Conflict(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusConflict, errs...)
}

// Response http.StatusUnprocessableEntity
func UnprocessableEntity(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error) {
	return WriteStatusAndError(req, resp, http.StatusUnprocessableEntity, errs...)
}
