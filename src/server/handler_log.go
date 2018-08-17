package server

import (
	"fmt"
	"net/http"

	"github.com/linkernetworks/vortex/src/container"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/net/http/query"
	"github.com/linkernetworks/vortex/src/web"
)

func getContainerLogsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	namespace := req.PathParameter("namespace")
	podID := req.PathParameter("pod")
	containerID := req.PathParameter("id")

	refTimestamp := req.QueryParameter("referenceTimestamp")
	if refTimestamp == "" {
		refTimestamp = container.NewestTimestamp
	}

	query := query.New(req.Request.URL.Query())

	refLineNum, err := query.Int("referenceLineNum", 0)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	tmp, ok := query.Str("previous")
	usePreviousLogs := tmp == "true"
	if ok == false {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	logFilePosition, ok := query.Str("logFilePosition")
	if ok == false {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	offsetFrom, err1 := query.Int("offsetFrom", 0)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	offsetTo, err2 := query.Int("offsetTo", 0)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	logSelector := container.DefaultSelection
	if err1 == nil && err2 == nil {
		logSelector = &container.Selection{
			ReferencePoint: container.LogLineId{
				LogTimestamp: container.LogTimestamp(refTimestamp),
				LineNum:      refLineNum,
			},
			OffsetFrom:      offsetFrom,
			OffsetTo:        offsetTo,
			LogFilePosition: logFilePosition,
		}
	}

	result, err := container.GetLogDetails(sp, namespace, podID, containerID, logSelector, usePreviousLogs)
	if err != nil {
		fmt.Errorf("failed to get the details of node: %v", err)
		return
	}
	resp.WriteHeaderAndEntity(http.StatusOK, result)
}
