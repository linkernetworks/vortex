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
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/net/http/query"
	"github.com/linkernetworks/vortex/src/server/backend"
	"github.com/linkernetworks/vortex/src/service"
	"github.com/linkernetworks/vortex/src/web"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func createServiceHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	userID, ok := req.Attribute("UserID").(string)
	if !ok {
		response.Unauthorized(req.Request, resp.ResponseWriter, fmt.Errorf("Unauthorized: User ID is empty"))
		return
	}

	s := entity.Service{}
	if err := req.ReadEntity(&s); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := sp.Validator.Struct(s); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	session := sp.Mongo.NewSession()
	session.C(entity.ServiceCollectionName).EnsureIndex(mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	})
	defer session.Close()

	// Check whether this name has been used
	s.ID = bson.NewObjectId()
	s.CreatedAt = timeutils.Now()
	s.OwnerID = bson.ObjectIdHex(userID)
	if err := service.CreateService(sp, &s); err != nil {
		if errors.IsAlreadyExists(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Service Name: %s already existed", s.Name))
		} else if errors.IsConflict(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Create setting has conflict: %v", err))
		} else if errors.IsInvalid(err) {
			response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("Create setting is invalid: %v", err))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}

	if err := session.Insert(entity.ServiceCollectionName, &s); err != nil {
		if mgo.IsDup(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Service Name: %s already existed", s.Name))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}

	// find owner in user entity
	s.CreatedBy, _ = backend.FindUserByID(session, s.OwnerID)
	resp.WriteHeaderAndEntity(http.StatusCreated, s)
}

func deleteServiceHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	id := req.PathParameter("id")

	session := sp.Mongo.NewSession()
	defer session.Close()

	s := entity.Service{}
	if err := session.FindOne(entity.ServiceCollectionName, bson.M{"_id": bson.ObjectIdHex(id)}, &s); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := service.DeleteService(sp, &s); err != nil {
		if errors.IsNotFound(err) {
			response.NotFound(req.Request, resp.ResponseWriter, err)
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}

	if err := session.Remove(entity.ServiceCollectionName, "_id", bson.ObjectIdHex(id)); err != nil {
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

func listServiceHandler(ctx *web.Context) {
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

	services := []entity.Service{}
	var c = session.C(entity.ServiceCollectionName)
	var q *mgo.Query

	selector := bson.M{}
	q = c.Find(selector).Sort("_id").Skip((page - 1) * pageSize).Limit(pageSize)

	if err := q.All(&services); err != nil {
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
	for _, service := range services {
		// find owner in user entity
		service.CreatedBy, _ = backend.FindUserByID(session, service.OwnerID)
	}

	count, err := session.Count(entity.ServiceCollectionName, bson.M{})
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))
	resp.AddHeader("X-Total-Count", strconv.Itoa(count))
	resp.AddHeader("X-Total-Pages", strconv.Itoa(totalPages))
	resp.WriteEntity(services)
}

func getServiceHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	id := req.PathParameter("id")

	session := sp.Mongo.NewSession()
	defer session.Close()
	c := session.C(entity.ServiceCollectionName)

	var service entity.Service
	if err := c.FindId(bson.ObjectIdHex(id)).One(&service); err != nil {
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
	service.CreatedBy, _ = backend.FindUserByID(session, service.OwnerID)
	resp.WriteEntity(service)
}

func uploadServiceYAMLHandler(ctx *web.Context) {
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

	serviceObj := obj.(*v1.Service)

	var servicePorts []entity.ServicePort
	for _, port := range serviceObj.Spec.Ports {
		servicePorts = append(servicePorts, entity.ServicePort{
			Name:       port.Name,
			Port:       port.Port,
			TargetPort: int(port.TargetPort.IntVal),
			NodePort:   port.NodePort,
		})
	}

	var serviceType string
	serviceType = string(serviceObj.Spec.Type)
	if serviceType == "" {
		serviceType = "ClusterIP"
	}

	d := entity.Service{
		ID:        bson.NewObjectId(),
		OwnerID:   bson.ObjectIdHex(userID),
		Name:      serviceObj.ObjectMeta.Name,
		Namespace: serviceObj.ObjectMeta.Namespace,
		Type:      serviceType,
		Selector:  serviceObj.Spec.Selector,
		Ports:     servicePorts,
	}

	if d.Namespace == "" {
		d.Namespace = "default"
	}

	session := sp.Mongo.NewSession()
	session.C(entity.ServiceCollectionName).EnsureIndex(mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	})
	defer session.Close()

	d.CreatedAt = timeutils.Now()
	_, err = sp.KubeCtl.CreateService(serviceObj, d.Namespace)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Service Name: %s already existed", d.Name))
		} else if errors.IsConflict(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Create setting has conflict: %v", err))
		} else if errors.IsInvalid(err) {
			response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("Create setting is invalid: %v", err))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}

	if err := session.Insert(entity.ServiceCollectionName, &d); err != nil {
		if mgo.IsDup(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Service Name: %s already existed", d.Name))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}
	// find owner in user entity
	d.CreatedBy, _ = backend.FindUserByID(session, d.OwnerID)
	resp.WriteHeaderAndEntity(http.StatusCreated, d)
}
