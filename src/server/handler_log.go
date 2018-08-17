package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/linkernetworks/vortex/src/container"
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

	refLineNum, err := strconv.Atoi(req.QueryParameter("referenceLineNum"))
	if err != nil {
		refLineNum = 0
	}
	usePreviousLogs := req.QueryParameter("previous") == "true"
	offsetFrom, err1 := strconv.Atoi(req.QueryParameter("offsetFrom"))
	offsetTo, err2 := strconv.Atoi(req.QueryParameter("offsetTo"))
	logFilePosition := req.QueryParameter("logFilePosition")

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
		fmt.Errorf("Failed to get the details of node: %v", err)
		return
	}
	resp.WriteHeaderAndEntity(http.StatusOK, result)
}
