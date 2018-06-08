package web

import (
	"net/http"

	"bitbucket.org/linkernetworks/vortex/src/serviceprovider"
	restful "github.com/emicklei/go-restful"
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
