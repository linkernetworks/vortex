package server

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/linkernetworks/utils/timeutils"
	"github.com/linkernetworks/vortex/src/deployment"
	"github.com/linkernetworks/vortex/src/entity"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/net/http/query"
	"github.com/linkernetworks/vortex/src/web"
	"k8s.io/apimachinery/pkg/api/errors"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func createAppHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	p := entity.Deployment{}
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
	defer session.Close()

	// Check whether this name has been used
	p.ID = bson.NewObjectId()
	p.CreatedAt = timeutils.Now()
	if err := deployment.CheckDeploymentParameter(sp, &p); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := deployment.CreateDeployment(sp, &p); err != nil {
		if errors.IsAlreadyExists(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Deployment Name: %s already existed", p.Name))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}
	if err := session.Insert(entity.DeploymentCollectionName, &p); err != nil {
		if mgo.IsDup(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Deployment Name: %s already existed", p.Name))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}
	resp.WriteHeaderAndEntity(http.StatusCreated, p)
}

func deleteAppHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	id := req.PathParameter("id")

	session := sp.Mongo.NewSession()
	defer session.Close()

	p := entity.Deployment{}
	if err := session.FindOne(entity.DeploymentCollectionName, bson.M{"_id": bson.ObjectIdHex(id)}, &p); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := deployment.DeleteDeployment(sp, &p); err != nil {
		if errors.IsNotFound(err) {
			response.NotFound(req.Request, resp.ResponseWriter, err)
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}

	if err := session.Remove(entity.DeploymentCollectionName, "_id", bson.ObjectIdHex(id)); err != nil {
		switch err {
		case mgo.ErrNotFound:
			response.NotFound(req.Request, resp.ResponseWriter, err)
			return
		default:
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
	}

	resp.WriteEntity(response.ActionResponse{
		Error:   false,
		Message: "Delete success",
	})
}

func listAppHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	var pageSize = 10
	query := query.New(req.Request.URL.Query())

	page, err := query.Int("page", 1)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}
	pageSize, err = query.Int("page_size", pageSize)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	session := sp.Mongo.NewSession()
	defer session.Close()

	deployments := []entity.Deployment{}
	var c = session.C(entity.DeploymentCollectionName)
	var q *mgo.Query

	selector := bson.M{}
	q = c.Find(selector).Sort("_id").Skip((page - 1) * pageSize).Limit(pageSize)

	if err := q.All(&deployments); err != nil {
		switch err {
		case mgo.ErrNotFound:
			response.NotFound(req.Request, resp.ResponseWriter, err)
			return
		default:
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
	}

	count, err := session.Count(entity.DeploymentCollectionName, bson.M{})
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))
	resp.AddHeader("X-Total-Count", strconv.Itoa(count))
	resp.AddHeader("X-Total-Pages", strconv.Itoa(totalPages))
	resp.WriteEntity(deployments)
}

func getAppHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	id := req.PathParameter("id")

	session := sp.Mongo.NewSession()
	defer session.Close()
	c := session.C(entity.DeploymentCollectionName)

	var deployment entity.Deployment
	if err := c.FindId(bson.ObjectIdHex(id)).One(&deployment); err != nil {
		switch err {
		case mgo.ErrNotFound:
			response.NotFound(req.Request, resp.ResponseWriter, err)
			return
		default:
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
	}
	resp.WriteEntity(deployment)
}
