package server

import (
	"github.com/emicklei/go-restful"
	"github.com/gorilla/mux"
	handler "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

func (a *App) AppRoute() *mux.Router {
	router := mux.NewRouter()

	container := restful.NewContainer()

	container.Filter(globalLogging)

	container.Add(newVersionService(a.ServiceProvider))
	container.Add(newNetworkService(a.ServiceProvider))
	container.Add(newStorageService(a.ServiceProvider))
	container.Add(newVolumeService(a.ServiceProvider))
	container.Add(newMonitoringService(a.ServiceProvider))

	router.PathPrefix("/v1/").Handler(container)
	return router
}

func newVersionService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/versions").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.GET("/").To(versionHandler(sp)))
	return webService
}

func newNetworkService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/networks").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.GET("/").To(handler.RESTfulServiceHandler(sp, listNetworkHandler)))
	webService.Route(webService.GET("/{id}").To(handler.RESTfulServiceHandler(sp, getNetworkHandler)))
	webService.Route(webService.POST("/").To(handler.RESTfulServiceHandler(sp, createNetworkHandler)))
	webService.Route(webService.DELETE("/{id}").To(handler.RESTfulServiceHandler(sp, deleteNetworkHandler)))
	return webService
}

func newStorageService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/storage").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.POST("/").To(handler.RESTfulServiceHandler(sp, createStorage)))
	webService.Route(webService.GET("/").To(handler.RESTfulServiceHandler(sp, listStorage)))
	webService.Route(webService.DELETE("/{id}").To(handler.RESTfulServiceHandler(sp, deleteStorage)))
	return webService
}

func newVolumeService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/volume").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.POST("/").To(handler.RESTfulServiceHandler(sp, createVolume)))
	webService.Route(webService.DELETE("/{id}").To(handler.RESTfulServiceHandler(sp, deleteVolume)))
	webService.Route(webService.GET("/").To(handler.RESTfulServiceHandler(sp, listVolume)))
	return webService
}

func newMonitoringService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/monitoring").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.GET("/query").To(handler.RESTfulServiceHandler(sp, queryMetrics)))
	// node
	webService.Route(webService.GET("/nodes").To(handler.RESTfulServiceHandler(sp, listNodeMetricsHandler)))
	webService.Route(webService.GET("/nodes/{id}").To(handler.RESTfulServiceHandler(sp, getNodeMetricsHandler)))
	// pod
	webService.Route(webService.GET("/pods").To(handler.RESTfulServiceHandler(sp, listPodMetricsHandler)))
	webService.Route(webService.GET("/pods/{id}").To(handler.RESTfulServiceHandler(sp, getPodMetricsHandler)))
	// container
	webService.Route(webService.GET("/containers").To(handler.RESTfulServiceHandler(sp, listContainerMetricsHandler)))
	webService.Route(webService.GET("/containers/{id}").To(handler.RESTfulServiceHandler(sp, getContainerMetricsHandler)))
	// service
	webService.Route(webService.GET("/services").To(handler.RESTfulServiceHandler(sp, listServiceMetricsHandler)))
	webService.Route(webService.GET("/services/{id}").To(handler.RESTfulServiceHandler(sp, getServiceMetricsHandler)))
	// deployment
	webService.Route(webService.GET("/deployments").To(handler.RESTfulServiceHandler(sp, listDeploymentMetricsHandler)))
	webService.Route(webService.GET("/deployments/{id}").To(handler.RESTfulServiceHandler(sp, getDeploymentMetricsHandler)))
	return webService
}
