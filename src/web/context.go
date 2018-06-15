package web

import (
	"net/http"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

type Context struct {
	ServiceProvider *serviceprovider.Container
	Request         *restful.Request
	Response        *restful.Response
}

type NativeContext struct {
	ServiceProvider *serviceprovider.Container
	Request         *http.Request
	Response        http.ResponseWriter
}
