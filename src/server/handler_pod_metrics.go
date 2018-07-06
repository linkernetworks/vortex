package server

import (
	"github.com/linkernetworks/vortex/src/entity"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/net/http/query"
	pc "github.com/linkernetworks/vortex/src/prometheuscontroller"
	"github.com/linkernetworks/vortex/src/web"
)

func listPodMetricsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	query := query.New(req.Request.URL.Query())
	replacer := map[string]string{}

	replacer["metric"] = "kube_pod_info"

	if node, ok := query.Str("node"); ok {
		replacer["node"] = node
	} else {
		replacer["node"] = ".*"
	}

	if namespace, ok := query.Str("namespace"); ok {
		replacer["namespace"] = namespace
	} else {
		replacer["namespace"] = ".*"
	}

	podList, err := pc.ListResource(sp, "pod", replacer)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	resp.WriteEntity(podList)
}

func getPodMetricsHandler(ctx *web.Context) {
	_, _, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	pod := entity.PodMetrics{}

	resp.WriteEntity(pod)
}
