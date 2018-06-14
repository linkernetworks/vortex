package responsetest

import (
	"net/http/httptest"
	"testing"

	response "bitbucket.org/linkernetworks/aurora/src/net/http"

	"encoding/json"
	"encoding/xml"

	"github.com/stretchr/testify/assert"
)

func AssertStatusEqual(t *testing.T, resp *httptest.ResponseRecorder, status int) {
	assert.Equal(t, status, resp.Code)
}

func AssertErrorMessage(t *testing.T, resp *httptest.ResponseRecorder, msg string) (err error) {
	var payload = response.ErrorPayload{}
	var out = resp.Body.Bytes()
	var contentType = resp.Header().Get("content-type")

	switch contentType {
	case "application/json", "text/json":
		err = json.Unmarshal(out, &payload)
	case "application/xml", "text/xml":
		err = xml.Unmarshal(out, &payload)
	}
	assert.Equal(t, msg, payload.Message)
	return err
}

func AssertError(t *testing.T, resp *httptest.ResponseRecorder) (err error) {
	var payload = response.ErrorPayload{}
	var out = resp.Body.Bytes()
	var contentType = resp.Header().Get("content-type")

	switch contentType {
	case "application/json", "text/json":
		err = json.Unmarshal(out, &payload)
	case "application/xml", "text/xml":
		err = xml.Unmarshal(out, &payload)
	}
	assert.NoError(t, err)
	return err
}
