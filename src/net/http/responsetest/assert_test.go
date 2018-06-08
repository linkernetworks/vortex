package responsetest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	response "bitbucket.org/linkernetworks/aurora/src/net/http"

	"github.com/stretchr/testify/assert"
)

func ExampleAssertError() {
	t := &testing.T{}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "http://here.com/v1/signin", nil)
	if err != nil {
		panic(err)
	}
	response.InternalServerError(request, recorder, errors.New("Failed to do something"))
	AssertError(t, recorder)
}

func ExampleAssertErrorMessage() {
	t := &testing.T{}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "http://here.com/v1/signin", nil)
	if err != nil {
		panic(err)
	}
	response.InternalServerError(request, recorder, errors.New("Failed to do something"))
	AssertErrorMessage(t, recorder, "Failed to do something")
}

func TestAssertError(t *testing.T) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "http://here.com/v1/signin", nil)
	assert.NoError(t, err)

	wl, err := response.InternalServerError(request, recorder, errors.New("Failed to do something"))
	assert.NoError(t, err)
	assert.True(t, wl > 0)
	AssertError(t, recorder)
}

func TestAssertStatusEqual(t *testing.T) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "http://here.com/v1/signin", nil)
	assert.NoError(t, err)

	wl, err := response.InternalServerError(request, recorder, errors.New("Failed to do something"))
	assert.NoError(t, err)
	assert.True(t, wl > 0)
	AssertStatusEqual(t, recorder, http.StatusInternalServerError)
}

func TestAssertErrorMessage(t *testing.T) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "http://here.com/v1/signin", nil)
	assert.NoError(t, err)

	msg := "Failed to do something"
	wl, err := response.InternalServerError(request, recorder, errors.New(msg))
	assert.NoError(t, err)
	assert.True(t, wl > 0)
	AssertErrorMessage(t, recorder, msg)
}
