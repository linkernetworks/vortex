package server

import (
	handler "bitbucket.org/linkernetworks/vortex/src/net/http"
	"bitbucket.org/linkernetworks/vortex/src/serviceprovider"
	"github.com/emicklei/go-restful"
	"github.com/gorilla/mux"
)

func (a *App) AppRoute() *mux.Router {
	router := mux.NewRouter()

	container := restful.NewContainer()

	container.Filter(globalLogging)

	container.Add(newVersionService(a.ServiceProvider))
	container.Add(newNetworkService(a.ServiceProvider))

	router.PathPrefix("/v1/").Handler(container)
	return router
}

func newVersionService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/versions").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.GET("/").To(VersionHandler(sp)))
	return webService
}

func newNetworkService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/networks").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.GET("/").To(handler.RESTfulServiceHandler(sp, ListNetworkHandler)))
	webService.Route(webService.GET("/{id}").To(handler.RESTfulServiceHandler(sp, GetNetworkHandler)))
	webService.Route(webService.POST("/").To(handler.RESTfulServiceHandler(sp, CreateNetworkHandler)))
	webService.Route(webService.PUT("/{id}").To(handler.RESTfulServiceHandler(sp, UpdateNetworkHandler)))
	webService.Route(webService.DELETE("/{id}").To(handler.RESTfulServiceHandler(sp, DeleteNetworkHandler)))
	return webService
}
