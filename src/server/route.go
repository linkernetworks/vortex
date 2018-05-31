package server

import (
	"bitbucket.org/linkernetworks/vortex/src/serviceprovider"
	"github.com/emicklei/go-restful"
	"github.com/gorilla/mux"
)

func (a *App) AppRoute() *mux.Router {
	router := mux.NewRouter()

	container := restful.NewContainer()

	container.Filter(globalLogging)

	container.Add(newVersionService(a.ServiceProvider))

	router.PathPrefix("/v1/").Handler(container)
	return router
}

func newVersionService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/versions").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.GET("/").To(VersionHandler(sp)))
	return webService
}
