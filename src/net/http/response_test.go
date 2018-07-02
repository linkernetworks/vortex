package http

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleWriteStatusAndError() {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "http://here.com/v1/signin", nil)
	if err != nil {
		panic(err)
	}
	WriteStatusAndError(request, recorder, http.StatusBadRequest, errors.New("bad request"))
}

func TestEncodeErrorPayload(t *testing.T) {
	testCases := []struct {
		cases       string
		contentType string
		expType     string
		expMessage  string
	}{
		{"json", "text/json", "text/json", `{"error":false,"message":""}`},
		{"xml", "text/xml", "text/xml", `<response><error>false</error><message></message></response>`},
		{"default", "", "application/json", `{"error":false,"message":""}`},
	}

	for _, tc := range testCases {
		t.Run(tc.cases, func(t *testing.T) {
			errPayload := ErrorPayload{}

			request, err := http.NewRequest("POST", "http://here.com/v1/signin", nil)
			assert.NoError(t, err)

			request.Header.Set("Content-Type", tc.contentType)
			out, cType, err := EncodeErrorPayload(request, errPayload)
			assert.Equal(t, tc.expType, cType)
			assert.NoError(t, err)
			assert.Equal(t, tc.expMessage, string(out[:len(out)]))
		})
	}
}

func TestNewErrorPayload(t *testing.T) {
	errs := []error{
		fmt.Errorf("Error One"),
		fmt.Errorf("Error Two"),
	}

	err := NewErrorPayload(errs[0], errs[1])
	assert.Equal(t, errs[0].Error(), err.Message)
	assert.Equal(t, errs[1].Error(), err.PreviousMessage)
}

func TestWriteStatusAndError(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest("POST", "http://here.com/v1/signin", nil)
	assert.NoError(t, err)
	wl, err := WriteStatusAndError(request, recorder, http.StatusForbidden, errors.New("Error one"))
	assert.NoError(t, err)
	assert.True(t, wl > 0)

	//Default type is application/json
	assert.Equal(t, http.StatusForbidden, recorder.Result().StatusCode)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
}

func TestSetStatus(t *testing.T) {
	testCases := []struct {
		cases      string
		handler    func(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error)
		statusCode int
	}{
		{"Forbidden", Forbidden, http.StatusForbidden},
		{"BadRequest", BadRequest, http.StatusBadRequest},
		{"OK", OK, http.StatusOK},
		{"NotFound", NotFound, http.StatusNotFound},
		{"Unauthorized", Unauthorized, http.StatusUnauthorized},
		{"InternalServerError", InternalServerError, http.StatusInternalServerError},
		{"Conflict", Conflict, http.StatusConflict},
		{"UnprocessableEntity", UnprocessableEntity, http.StatusUnprocessableEntity},
	}

	for _, tc := range testCases {
		t.Run(tc.cases, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest("POST", "http://here.com/v1/signin", nil)
			assert.NoError(t, err)
			wl, err := tc.handler(request, recorder, errors.New("Failed to do something"))
			assert.NoError(t, err)
			assert.True(t, wl > 0)
			assert.Equal(t, tc.statusCode, recorder.Result().StatusCode)

		})
	}
}
