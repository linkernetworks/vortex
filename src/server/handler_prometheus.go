package server

import (
	"fmt"
	"net/http"
	"time"

	restful "github.com/emicklei/go-restful"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/net/http/query"
	"github.com/linkernetworks/vortex/src/web"
	"golang.org/x/net/context"
)

func queryMetrics(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	query := query.New(req.Request.URL.Query())

	query_str := ""
	if q, ok := query.Str("query"); ok {
		query_str = q
	}

	api := sp.Prometheus.API

	testTime := time.Now()
	result, err := api.Query(context.Background(), query_str, testTime)

	if result == nil {
		response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("%v: %v", result, err))
	}

	resp.WriteJson(map[string]interface{}{
		"status":  http.StatusOK,
		"results": result,
	}, restful.MIME_JSON)
}
