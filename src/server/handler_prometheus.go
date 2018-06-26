package server

import (
	"net/http"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/vortex/src/net/http/query"
	"github.com/linkernetworks/vortex/src/web"
	"golang.org/x/net/context"
)

func queryMetrics(ctx *web.Context) {
	as, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	query := query.New(req.Request.URL.Query())

	query_str := ""
	if q, ok := query.Str("query"); ok {
		query_str = q
	}

	api := as.Prometheus.API

	testTime := time.Now()
	result, _ := api.Query(context.Background(), query_str, testTime)

	resp.WriteJson(map[string]interface{}{
		"status":  http.StatusOK,
		"results": result,
	}, restful.MIME_JSON)
}
