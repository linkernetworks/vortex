package server

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/linkernetworks/utils/timeutils"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/namespace"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/net/http/query"
	"github.com/linkernetworks/vortex/src/web"
	"k8s.io/apimachinery/pkg/api/errors"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func createNamespaceHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	n := entity.Namespace{}
	if err := req.ReadEntity(&n); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := sp.Validator.Struct(n); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	session := sp.Mongo.NewSession()
	session.C(entity.NamespaceCollectionName).EnsureIndex(mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	})
	defer session.Close()

	// Check whether this name has been used
	n.ID = bson.NewObjectId()
	n.CreatedAt = timeutils.Now()

	if err := namespace.CreateNamespace(sp, &n); err != nil {
		if errors.IsAlreadyExists(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Namespace: %s already existed", n.Name))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}
	if err := session.Insert(entity.NamespaceCollectionName, &n); err != nil {
		if mgo.IsDup(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Namespace: %s already existed", n.Name))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}
	resp.WriteHeaderAndEntity(http.StatusCreated, n)
}

func deleteNamespaceHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	id := req.PathParameter("id")

	session := sp.Mongo.NewSession()
	defer session.Close()

	n := entity.Namespace{}
	if err := session.FindOne(entity.NamespaceCollectionName, bson.M{"_id": bson.ObjectIdHex(id)}, &n); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := namespace.DeleteNamespace(sp, &n); err != nil {
		if errors.IsNotFound(err) {
			response.NotFound(req.Request, resp.ResponseWriter, err)
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}

	if err := session.Remove(entity.NamespaceCollectionName, "_id", bson.ObjectIdHex(id)); err != nil {
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

func listNamespaceHandler(ctx *web.Context) {
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

	namespaces := []entity.Namespace{}
	var c = session.C(entity.NamespaceCollectionName)
	var q *mgo.Query

	selector := bson.M{}
	q = c.Find(selector).Sort("_id").Skip((page - 1) * pageSize).Limit(pageSize)

	if err := q.All(&namespaces); err != nil {
		switch err {
		case mgo.ErrNotFound:
			response.NotFound(req.Request, resp.ResponseWriter, err)
			return
		default:
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
	}

	count, err := session.Count(entity.NamespaceCollectionName, bson.M{})
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))
	resp.AddHeader("X-Total-Count", strconv.Itoa(count))
	resp.AddHeader("X-Total-Pages", strconv.Itoa(totalPages))
	resp.WriteEntity(namespaces)
}

func getNamespaceHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	id := req.PathParameter("id")

	session := sp.Mongo.NewSession()
	defer session.Close()
	c := session.C(entity.NamespaceCollectionName)

	var namespace entity.Namespace
	if err := c.FindId(bson.ObjectIdHex(id)).One(&namespace); err != nil {
		switch err {
		case mgo.ErrNotFound:
			response.NotFound(req.Request, resp.ResponseWriter, err)
			return
		default:
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
	}
	resp.WriteEntity(namespace)
}
