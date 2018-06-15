package server

import (
	"github.com/emicklei/go-restful"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

func VersionHandler(sp *serviceprovider.Container) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
	}
}
