package server

import (
	"net/http"

	"github.com/linkernetworks/vortex/src/entity"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/web"
)

func registryBasicAuthHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	credential := entity.RegistryBasicAuthCredential{}
	if err := req.ReadEntity(&credential); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := sp.Validator.Struct(credential); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	client := &http.Client{}
	registryReq, err := http.NewRequest("GET", sp.Config.Registry.URL+"/v2", nil)
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	registryReq.SetBasicAuth(credential.Username, credential.Password)
	registryResp, err := client.Do(registryReq)
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	resp.WriteHeaderAndEntity(registryResp.StatusCode, credential)
}
