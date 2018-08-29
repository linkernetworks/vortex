package server

import (
	"fmt"
	"net/http"

	"github.com/linkernetworks/utils/timeutils"
	"github.com/linkernetworks/vortex/src/deployment"
	"github.com/linkernetworks/vortex/src/entity"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/server/backend"
	"github.com/linkernetworks/vortex/src/service"
	"github.com/linkernetworks/vortex/src/web"
	"k8s.io/apimachinery/pkg/api/errors"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func createAppHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	userID, ok := req.Attribute("UserID").(string)
	if !ok {
		response.Unauthorized(req.Request, resp.ResponseWriter, fmt.Errorf("Unauthorized: User ID is empty"))
		return
	}

	p := entity.Application{}
	if err := req.ReadEntity(&p); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := sp.Validator.Struct(p); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	session := sp.Mongo.NewSession()
	session.C(entity.DeploymentCollectionName).EnsureIndex(mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	})
	session.C(entity.ServiceCollectionName).EnsureIndex(mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	})
	defer session.Close()

	p.Deployment.ID = bson.NewObjectId()
	p.Service.ID = bson.NewObjectId()

	p.Deployment.OwnerID = bson.ObjectIdHex(userID)
	p.Service.OwnerID = bson.ObjectIdHex(userID)

	time := timeutils.Now()
	p.Deployment.CreatedAt = time
	p.Service.CreatedAt = time
	if err := deployment.CheckDeploymentParameter(sp, &p.Deployment); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	//Bond to same namespace
	p.Service.Namespace = p.Deployment.Namespace

	//We use the application label for deployment and service
	p.Service.Selector[deployment.DefaultLabel] = p.Deployment.Name
	if err := deployment.CreateDeployment(sp, &p.Deployment); err != nil {
		if errors.IsAlreadyExists(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Deployment Name: %s already existed", p.Deployment.Name))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}

	if err := service.CreateService(sp, &p.Service); err != nil {
		if errors.IsAlreadyExists(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Service Name: %s already existed", p.Service.Name))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}

	p.Deployment.OwnerID = bson.ObjectIdHex(userID)
	if err := session.Insert(entity.DeploymentCollectionName, &p.Deployment); err != nil {
		if mgo.IsDup(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Deployment Name: %s already existed", p.Deployment.Name))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}

	p.Service.OwnerID = bson.ObjectIdHex(userID)
	if err := session.Insert(entity.ServiceCollectionName, &p.Service); err != nil {
		if mgo.IsDup(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Service Name: %s already existed", p.Service.Name))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}
	// find owner in user entity
	p.Service.CreatedBy, _ = backend.FindUserByID(session, p.Service.OwnerID)
	p.Deployment.CreatedBy, _ = backend.FindUserByID(session, p.Deployment.OwnerID)
	resp.WriteHeaderAndEntity(http.StatusCreated, p)
}
