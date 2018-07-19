package server

import (
	"github.com/linkernetworks/vortex/src/version"
	"github.com/linkernetworks/vortex/src/web"
)

func versionHandler(ctx *web.Context) {
	_, _, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	resp.WriteEntity(ActionResponse{
		Error:   false,
		Message: version.GetVersion(),
	})
}
