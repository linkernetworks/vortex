package server

import (
	"github.com/emicklei/go-restful"
	"github.com/linkernetworks/logger"
)

// SessionKey will be the cookie name defined in the http header
const SessionKey = "ses"

func globalLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	logger.Infof("%s %s", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
}
