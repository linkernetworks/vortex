package server

import (
	"fmt"
	"net/http"

	restful "github.com/emicklei/go-restful"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/web"
	"github.com/prometheus/common/model"
)

func listNodeHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	result, err := queryFromPrometheus(sp, "sum by (node)(kube_node_info)")

	if result == nil {
		response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("%v: %v", result, err))
	}

	nodeList := []model.LabelValue{}

	switch {
	case result.Type() == model.ValVector:
		result := result.(model.Vector)
		for _, node := range result {
			nodeList = append(nodeList, node.Metric["node"])
		}
	}

	resp.WriteJson(map[string]interface{}{
		"status":  http.StatusOK,
		"results": nodeList,
	}, restful.MIME_JSON)
}

func getNodeHandler(ctx *web.Context) {
	//sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

}
