package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"net/http"

	"github.com/linkernetworks/vortex/src/web"
)

func TestCompositeServiceHandler(t *testing.T) {
	data := 0
	var handler = func(*web.NativeContext) {
		data++
	}

	req, err := http.NewRequest("POST", "http://here.com/v1/signin", nil)
	assert.NoError(t, err)

	routeHandler := CompositeServiceHandler(nil, handler)
	assert.Equal(t, 0, data)
	routeHandler(nil, req)
	assert.Equal(t, 1, data)
}

func TestRESTfulServiceHandler(t *testing.T) {
	data := 0
	var handler = func(*web.Context) {
		data++
	}

	routeHandler := RESTfulServiceHandler(nil, handler)
	assert.Equal(t, 0, data)
	routeHandler(nil, nil)
	assert.Equal(t, 1, data)
}
