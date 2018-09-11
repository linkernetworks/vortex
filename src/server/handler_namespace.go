package server

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"

	"github.com/linkernetworks/utils/timeutils"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/kubernetes"
	"github.com/linkernetworks/vortex/src/namespace"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/net/http/query"
	"github.com/linkernetworks/vortex/src/server/backend"
	"github.com/linkernetworks/vortex/src/web"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func createNamespaceHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	userID, ok := req.Attribute("UserID").(string)
	if !ok {
		response.Unauthorized(req.Request, resp.ResponseWriter, fmt.Errorf("Unauthorized: User ID is empty"))
		return
	}

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

	n.OwnerID = bson.ObjectIdHex(userID)
	if err := session.Insert(entity.NamespaceCollectionName, &n); err != nil {
		if mgo.IsDup(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Namespace: %s already existed", n.Name))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}
	n.CreatedBy, _ = backend.FindUserByID(session, n.OwnerID)
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

	// insert users entity
	for _, namespace := range namespaces {
		// find owner in user entity
		namespace.CreatedBy, _ = backend.FindUserByID(session, namespace.OwnerID)
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
	// find owner in user entity
	namespace.CreatedBy, _ = backend.FindUserByID(session, namespace.OwnerID)
	resp.WriteEntity(namespace)
}

func uploadNamespaceYAMLHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	userID, ok := req.Attribute("UserID").(string)
	if !ok {
		response.Unauthorized(req.Request, resp.ResponseWriter, fmt.Errorf("Unauthorized: User ID is empty"))
		return
	}

	if err := req.Request.ParseMultipartForm(_24K); nil != err {
		response.InternalServerError(req.Request, resp.ResponseWriter, fmt.Errorf("Failed to read multipart form: %s", err.Error()))
		return
	}

	infile, _, err := req.Request.FormFile("file")
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("Error parsing uploaded file %v", err))
		return
	}

	content, err := ioutil.ReadAll(infile)
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, fmt.Errorf("Failed to read data: %s", err.Error()))
		return
	}

	if len(content) == 0 {
		response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("Empty content"))
		return
	}

	obj, err := kubernetes.ParseK8SYAML(content)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	namespaceObj := obj.(*v1.Namespace)

	d := entity.Namespace{
		ID:      bson.NewObjectId(),
		OwnerID: bson.ObjectIdHex(userID),
		Name:    namespaceObj.ObjectMeta.Name,
	}

	session := sp.Mongo.NewSession()
	session.C(entity.NamespaceCollectionName).EnsureIndex(mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	})
	defer session.Close()

	d.CreatedAt = timeutils.Now()
	_, err = sp.KubeCtl.CreateNamespace(namespaceObj)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Namespace Name: %s already existed", d.Name))
		} else if errors.IsConflict(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Create setting has conflict: %v", err))
		} else if errors.IsInvalid(err) {
			response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("Create setting is invalid: %v", err))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}

	if err := session.Insert(entity.NamespaceCollectionName, &d); err != nil {
		if mgo.IsDup(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Namespace Name: %s already existed", d.Name))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}
	// find owner in user entity
	d.CreatedBy, _ = backend.FindUserByID(session, d.OwnerID)
	resp.WriteHeaderAndEntity(http.StatusCreated, d)
}
