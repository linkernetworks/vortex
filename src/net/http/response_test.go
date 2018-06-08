package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"
	"encoding/xml"

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

func ExampleDefaultEncoding() {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "http://here.com/v1/signin", nil)
	if err != nil {
		panic(err)
	}
	InternalServerError(request, recorder, errors.New("Failed to do something"))
}

func ExampleJsonEncoding() {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "http://here.com/v1/signin", nil)
	if err != nil {
		panic(err)
	}

	request.Header.Set("content-type", "application/json")
	InternalServerError(request, recorder, errors.New("Failed to do something"))
}

func TestDefaultErrorXmlEncode(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest("POST", "http://here.com/v1/signin", nil)
	request.Header.Set("content-type", "application/xml")
	assert.NoError(t, err)

	msg := "Failed to do something"

	wl, err := InternalServerError(request, recorder, errors.New(msg))
	assert.NoError(t, err)
	assert.True(t, wl > 0)

	contentType := recorder.Header().Get("content-type")
	assert.Equal(t, "application/xml", contentType)

	payload := ErrorPayload{}
	out := recorder.Body.Bytes()
	err = xml.Unmarshal(out, &payload)
	assert.NoError(t, err)

	t.Logf("XML response: %s", out)

	assert.True(t, payload.Error)
	assert.Len(t, payload.Message, len(msg))
}

func TestDefaultErrorJsonEncode(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest("POST", "http://here.com/v1/signin", nil)
	assert.NoError(t, err)
	wl, err := InternalServerError(request, recorder, errors.New("Failed to do something"))
	assert.NoError(t, err)
	assert.True(t, wl > 0)

	msg := "Failed to do something"

	contentType := recorder.Header().Get("content-type")
	assert.Equal(t, "application/json", contentType)

	payload := ErrorPayload{}
	out := recorder.Body.Bytes()
	err = json.Unmarshal(out, &payload)
	assert.NoError(t, err)

	assert.True(t, payload.Error)
	assert.Len(t, payload.Message, len(msg))
}
