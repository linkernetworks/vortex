package http

import (
	"net/http"

	"bitbucket.org/linkernetworks/vortex/src/serviceprovider"
	"bitbucket.org/linkernetworks/vortex/src/web"
	"github.com/emicklei/go-restful"
	"github.com/linkernetworks/logger"
)

type NativeContextHandler func(*web.NativeContext)

// CompositeServiceProvider apply mongo client to HandlerFunc
func CompositeServiceHandler(sp *serviceprovider.Container, handler NativeContextHandler) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		logger.Infoln(req.Method, req.URL)
		ctx := web.NativeContext{sp, req, resp}
		handler(&ctx)
	}
}

type RESTfulContextHandler func(*web.Context)

func RESTfulServiceHandler(sp *serviceprovider.Container, handler RESTfulContextHandler) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		ctx := web.Context{sp, req, resp}
		handler(&ctx)
	}
}
