package http

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/linkernetworks/vortex/src/web"
)

// The interface for native http handler(http.Request and http.ResponseWriter)
type NativeContextHandler func(*web.NativeContext)

// CompositeServiceProvider apply mongo client to HandlerFunc
func CompositeServiceHandler(sp *serviceprovider.Container, handler NativeContextHandler) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		logger.Infoln(req.Method, req.URL)
		ctx := web.NativeContext{sp, req, resp}
		handler(&ctx)
	}
}

// The interface for restfuul handler(restful.Request,restful.Response)
type RESTfulContextHandler func(*web.Context)

// The wrapper to combine the RESTfulContextHandler with our serviceprovider object
func RESTfulServiceHandler(sp *serviceprovider.Container, handler RESTfulContextHandler) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		ctx := web.Context{sp, req, resp}
		handler(&ctx)
	}
}
