package http

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/linkernetworks/vortex/src/web"
)

// NativeContextHandler is the interface for native http handler(http.Request and http.ResponseWriter)
type NativeContextHandler func(*web.NativeContext)

// CompositeServiceHandler apply mongo client to HandlerFunc
func CompositeServiceHandler(sp *serviceprovider.Container, handler NativeContextHandler) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		logger.Infoln(req.Method, req.URL)
		ctx := web.NativeContext{
			ServiceProvider: sp,
			Request:         req,
			Response:        resp,
		}
		handler(&ctx)
	}
}

// RESTfulContextHandler is the interface for restfuul handler(restful.Request,restful.Response)
type RESTfulContextHandler func(*web.Context)

// RESTfulServiceHandler is the wrapper to combine the RESTfulContextHandler with our serviceprovider object
func RESTfulServiceHandler(sp *serviceprovider.Container, handler RESTfulContextHandler) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		ctx := web.Context{
			ServiceProvider: sp,
			Request:         req,
			Response:        resp,
		}
		handler(&ctx)
	}
}
