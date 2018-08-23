package server

import (
	"github.com/emicklei/go-restful"
	"github.com/gorilla/mux"
	handler "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

// AppRoute will add router
func (a *App) AppRoute() *mux.Router {
	router := mux.NewRouter()

	container := restful.NewContainer()

	container.Filter(globalLogging)

	container.Add(newVersionService(a.ServiceProvider))
	container.Add(newRegistryService(a.ServiceProvider))
	container.Add(newUserService(a.ServiceProvider))
	container.Add(newNetworkService(a.ServiceProvider))
	container.Add(newStorageService(a.ServiceProvider))
	container.Add(newVolumeService(a.ServiceProvider))
	container.Add(newContainerService(a.ServiceProvider))
	container.Add(newPodService(a.ServiceProvider))
	container.Add(newDeploymentService(a.ServiceProvider))
	container.Add(newServiceService(a.ServiceProvider))
	container.Add(newNamespaceService(a.ServiceProvider))
	container.Add(newMonitoringService(a.ServiceProvider))
	container.Add(newAppService(a.ServiceProvider))

	router.PathPrefix("/v1/").Handler(container)
	return router
}

func newVersionService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/version").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.GET("/").To(handler.RESTfulServiceHandler(sp, versionHandler)))
	return webService
}

func newRegistryService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/registry").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	//	webService.Filter(validateTokenMiddleware)
	webService.Route(webService.POST("/auth").To(handler.RESTfulServiceHandler(sp, registryBasicAuthHandler)))
	return webService
}

func newUserService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/users").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.POST("/").To(handler.RESTfulServiceHandler(sp, createUserHandler)))
	webService.Route(webService.DELETE("/{id}").To(handler.RESTfulServiceHandler(sp, deleteUserHandler)))
	webService.Route(webService.GET("/").To(handler.RESTfulServiceHandler(sp, listUserHandler)))
	webService.Route(webService.GET("/{id}").To(handler.RESTfulServiceHandler(sp, getUserHandler)))
	// Authenticate handlers Sign Up / Sign In
	webService.Route(webService.POST("/signup").To(handler.RESTfulServiceHandler(sp, signUpUserHandler)))
	webService.Route(webService.POST("/signin").To(handler.RESTfulServiceHandler(sp, signInUserHandler)))
	return webService
}

func newNetworkService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/networks").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Filter(validateTokenMiddleware)
	webService.Route(webService.GET("/").To(handler.RESTfulServiceHandler(sp, listNetworkHandler)))
	webService.Route(webService.GET("/{id}").To(handler.RESTfulServiceHandler(sp, getNetworkHandler)))
	webService.Route(webService.GET("/status/{id}").To(handler.RESTfulServiceHandler(sp, getNetworkStatusHandler)))
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
	webService.Filter(validateTokenMiddleware)
	webService.Route(webService.POST("/").To(handler.RESTfulServiceHandler(sp, createVolumeHandler)))
	webService.Route(webService.DELETE("/{id}").To(handler.RESTfulServiceHandler(sp, deleteVolumeHandler)))
	webService.Route(webService.GET("/").To(handler.RESTfulServiceHandler(sp, listVolumeHandler)))
	return webService
}

func newContainerService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/containers").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.GET("/logs/{namespace}/{pod}/{container}").To(handler.RESTfulServiceHandler(sp, getContainerLogsHandler)))
	webService.Route(webService.GET("/logs/file/{namespace}/{pod}/{container}").To(handler.RESTfulServiceHandler(sp, getContainerLogFileHandler)))
	return webService
}

func newPodService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/pods").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Filter(validateTokenMiddleware)
	webService.Route(webService.POST("/").To(handler.RESTfulServiceHandler(sp, createPodHandler)))
	webService.Route(webService.DELETE("/{id}").To(handler.RESTfulServiceHandler(sp, deletePodHandler)))
	webService.Route(webService.GET("/").To(handler.RESTfulServiceHandler(sp, listPodHandler)))
	webService.Route(webService.GET("/{id}").To(handler.RESTfulServiceHandler(sp, getPodHandler)))
	return webService
}

func newDeploymentService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/deployments").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.POST("/").To(handler.RESTfulServiceHandler(sp, createDeploymentHandler)))
	webService.Route(webService.DELETE("/{id}").To(handler.RESTfulServiceHandler(sp, deleteDeploymentHandler)))
	webService.Route(webService.GET("/").To(handler.RESTfulServiceHandler(sp, listDeploymentHandler)))
	webService.Route(webService.GET("/{id}").To(handler.RESTfulServiceHandler(sp, getDeploymentHandler)))
	return webService
}

func newAppService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/apps").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.POST("/").To(handler.RESTfulServiceHandler(sp, createAppHandler)))
	return webService
}

func newServiceService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/services").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Filter(validateTokenMiddleware)
	webService.Route(webService.POST("/").To(handler.RESTfulServiceHandler(sp, createServiceHandler)))
	webService.Route(webService.DELETE("/{id}").To(handler.RESTfulServiceHandler(sp, deleteServiceHandler)))
	webService.Route(webService.GET("/").To(handler.RESTfulServiceHandler(sp, listServiceHandler)))
	webService.Route(webService.GET("/{id}").To(handler.RESTfulServiceHandler(sp, getServiceHandler)))
	return webService
}

func newNamespaceService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/namespaces").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.POST("/").To(handler.RESTfulServiceHandler(sp, createNamespaceHandler)))
	webService.Route(webService.DELETE("/{id}").To(handler.RESTfulServiceHandler(sp, deleteNamespaceHandler)))
	webService.Route(webService.GET("/").To(handler.RESTfulServiceHandler(sp, listNamespaceHandler)))
	webService.Route(webService.GET("/{id}").To(handler.RESTfulServiceHandler(sp, getNamespaceHandler)))
	return webService
}

func newMonitoringService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/monitoring").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	// node
	webService.Route(webService.GET("/nodes").To(handler.RESTfulServiceHandler(sp, listNodeMetricsHandler)))
	webService.Route(webService.GET("/nodes/{id}").To(handler.RESTfulServiceHandler(sp, getNodeMetricsHandler)))
	webService.Route(webService.GET("/nodes/{id}/nics").To(handler.RESTfulServiceHandler(sp, listNodeNicsMetricsHandler)))
	// pod
	webService.Route(webService.GET("/pods").To(handler.RESTfulServiceHandler(sp, listPodMetricsHandler)))
	webService.Route(webService.GET("/pods/{id}").To(handler.RESTfulServiceHandler(sp, getPodMetricsHandler)))
	// container
	webService.Route(webService.GET("/containers").To(handler.RESTfulServiceHandler(sp, listContainerMetricsHandler)))
	webService.Route(webService.GET("/containers/{id}").To(handler.RESTfulServiceHandler(sp, getContainerMetricsHandler)))
	// service
	webService.Route(webService.GET("/services").To(handler.RESTfulServiceHandler(sp, listServiceMetricsHandler)))
	webService.Route(webService.GET("/services/{id}").To(handler.RESTfulServiceHandler(sp, getServiceMetricsHandler)))
	// controller
	webService.Route(webService.GET("/controllers").To(handler.RESTfulServiceHandler(sp, listControllerMetricsHandler)))
	webService.Route(webService.GET("/controllers/{id}").To(handler.RESTfulServiceHandler(sp, getControllerMetricsHandler)))
	return webService
}
