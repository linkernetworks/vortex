package server

import (
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/net/http/query"
	pc "github.com/linkernetworks/vortex/src/prometheuscontroller"
	"github.com/linkernetworks/vortex/src/web"
	"github.com/prometheus/common/model"
)

func listServiceMetricsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	query := query.New(req.Request.URL.Query())
	expression := pc.Expression{}
	expression.Metrics = []string{"kube_service_info"}
	expression.QueryLabels = map[string]string{}
	expression.TargetLabels = []model.LabelName{"service"}

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

	containerList, err := pc.ListResource(sp, expression)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	resp.WriteEntity(containerList)
}

func getServiceMetricsHandler(ctx *web.Context) {
	// _, _, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	// id := req.PathParameter("id")

	// pod := entity.PodMetrics{}

	// resp.WriteEntity()
}
