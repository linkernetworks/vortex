package server

import (
	"bitbucket.org/linkernetworks/vortex/src/serviceprovider"
	"github.com/emicklei/go-restful"
)

func VersionHandler(sp *serviceprovider.Container) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
	}
}
