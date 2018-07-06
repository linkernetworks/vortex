package server

import (
	"github.com/linkernetworks/vortex/src/entity"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/net/http/query"
	pc "github.com/linkernetworks/vortex/src/prometheuscontroller"
	"github.com/linkernetworks/vortex/src/web"
)

func listServiceMetricsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	query := query.New(req.Request.URL.Query())
	expression := pc.Expression{}
	expression.Metrics = []string{"kube_service_info"}
	expression.QueryLabels = map[string]string{}

	if node, ok := query.Str("node"); ok {
		expression.QueryLabels["node"] = node
	} else {
		expression.QueryLabels["node"] = ".*"
	}

	if namespace, ok := query.Str("namespace"); ok {
		expression.QueryLabels["namespace"] = namespace
	} else {
		expression.QueryLabels["namespace"] = ".*"
	}

	containerList, err := pc.ListResource(sp, "service", expression)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	resp.WriteEntity(containerList)
}

func getServiceMetricsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	service := entity.ServiceMetrics{}
	id := req.PathParameter("id")

	service, err := pc.GetService(sp, id)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	resp.WriteEntity(service)
}
