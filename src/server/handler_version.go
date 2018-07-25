package server

import (
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/version"
	"github.com/linkernetworks/vortex/src/web"
)

func versionHandler(ctx *web.Context) {
	_, _, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	resp.WriteEntity(response.ActionResponse{
		Error:   false,
		Message: version.GetVersion(),
	})
}
