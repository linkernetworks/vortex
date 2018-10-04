package server

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/linkernetworks/utils/timeutils"
	"github.com/linkernetworks/vortex/src/deployment"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/kubernetes"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/net/http/query"
	"github.com/linkernetworks/vortex/src/server/backend"
	"github.com/linkernetworks/vortex/src/web"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const _24K = (1 << 10) * 24

func createDeploymentHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	userID, ok := req.Attribute("UserID").(string)
	if !ok {
		response.Unauthorized(req.Request, resp.ResponseWriter, fmt.Errorf("Unauthorized: User ID is empty"))
		return
	}

	p := entity.Deployment{}
	if err := req.ReadEntity(&p); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	// autoscale resource request setting prerequisite check
	for _, container := range p.Containers {
		if container.ResourceRequestCPU != 0 {
			p.AutoscalerInfo.IsCapableAutoscaleResources = append(p.AutoscalerInfo.IsCapableAutoscaleResources, "cpu")
		}
		if container.ResourceRequestMemory != 0 {
			p.AutoscalerInfo.IsCapableAutoscaleResources = append(p.AutoscalerInfo.IsCapableAutoscaleResources, "memory")
		}
	}

	if p.NetworkType == entity.DeploymentCustomNetwork {
		// clear the slice
		p.AutoscalerInfo.IsCapableAutoscaleResources = nil
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

	p.OwnerID = bson.ObjectIdHex(userID)
	// find owner in user entity
	ownerUser, _ := backend.FindUserByID(session, p.OwnerID)

	var account, domain string
	components := strings.Split(ownerUser.LoginCredential.Username, "@")
	account, domain = components[0], components[1]

	// append label with owner email
	p.Labels[deployment.NotificationEmailAccount] = account
	p.Labels[deployment.NotificationEmailDomain] = domain

	if err := deployment.CreateDeployment(sp, &p); err != nil {
		if errors.IsAlreadyExists(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Deployment Name: %s already existed", p.Name))
		} else if errors.IsConflict(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Create setting has conflict: %v", err))
		} else if errors.IsInvalid(err) {
			response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("Create setting is invalid: %v", err))
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
	p.CreatedBy = ownerUser
	resp.WriteHeaderAndEntity(http.StatusCreated, p)
}

func deleteDeploymentHandler(ctx *web.Context) {
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

	autoscaler := entity.AutoscalerInfo{
		Namespace:          p.Namespace,
		ScaleTargetRefName: p.Name,
	}
	// Autoscaler should be clean if there is any and ignore delete autoscaler error
	deployment.DeleteAutoscaler(sp, autoscaler)
	resp.WriteEntity(response.ActionResponse{
		Error:   false,
		Message: "Delete success",
	})
}

func listDeploymentHandler(ctx *web.Context) {
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

	// insert users entity
	for _, deployment := range deployments {
		// find owner in user entity
		deployment.CreatedBy, _ = backend.FindUserByID(session, deployment.OwnerID)
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

func getDeploymentHandler(ctx *web.Context) {
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
	deployment.CreatedBy, _ = backend.FindUserByID(session, deployment.OwnerID)
	resp.WriteEntity(deployment)
}

func uploadDeploymentYAMLHandler(ctx *web.Context) {
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

	deploymentObj, ok := obj.(*appv1.Deployment)
	if !ok {
		response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("The YAML file is not for creating deployment"))
		return
	}

	d := entity.Deployment{
		ID:        bson.NewObjectId(),
		OwnerID:   bson.ObjectIdHex(userID),
		Name:      deploymentObj.ObjectMeta.Name,
		Namespace: deploymentObj.ObjectMeta.Namespace,
		Replicas:  *deploymentObj.Spec.Replicas,
	}

	if d.Namespace == "" {
		d.Namespace = "default"
	}

	// autoscale resource request setting prerequisite check
	for _, container := range deploymentObj.Spec.Template.Spec.Containers {
		// check if map contains cpu key
		if _, ok := container.Resources.Requests[corev1.ResourceCPU]; ok {
			d.AutoscalerInfo.IsCapableAutoscaleResources = append(d.AutoscalerInfo.IsCapableAutoscaleResources, "cpu")
		}
		// check if map contains memory key
		if _, ok := container.Resources.Requests[corev1.ResourceMemory]; ok {
			d.AutoscalerInfo.IsCapableAutoscaleResources = append(d.AutoscalerInfo.IsCapableAutoscaleResources, "memory")
		}
	}

	session := sp.Mongo.NewSession()
	session.C(entity.DeploymentCollectionName).EnsureIndex(mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	})
	defer session.Close()

	d.CreatedAt = timeutils.Now()

	// find owner in user entity
	ownerUser, _ := backend.FindUserByID(session, d.OwnerID)

	var account, domain string
	components := strings.Split(ownerUser.LoginCredential.Username, "@")
	account, domain = components[0], components[1]

	// append label with owner email
	deploymentObj.ObjectMeta.Labels[deployment.NotificationEmailAccount] = account
	deploymentObj.ObjectMeta.Labels[deployment.NotificationEmailDomain] = domain

	_, err = sp.KubeCtl.CreateDeployment(deploymentObj, d.Namespace)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Deployment Name: %s already existed", d.Name))
		} else if errors.IsConflict(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Create setting has conflict: %v", err))
		} else if errors.IsInvalid(err) {
			response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("Create setting is invalid: %v", err))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}

	if err := session.Insert(entity.DeploymentCollectionName, &d); err != nil {
		if mgo.IsDup(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Deployment Name: %s already existed", d.Name))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}
	// find owner in user entity
	d.CreatedBy = ownerUser
	resp.WriteHeaderAndEntity(http.StatusCreated, d)
}

func updateAutoscalerHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	enableAutoscaler := req.QueryParameter("enable") == "true"

	session := sp.Mongo.NewSession()
	defer session.Close()

	autoscalerInfo := entity.AutoscalerInfo{}
	if err := req.ReadEntity(&autoscalerInfo); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}
	// force set autoscaler to deployment name
	autoscalerInfo.Name = autoscalerInfo.ScaleTargetRefName
	if err := sp.Validator.Struct(autoscalerInfo); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	if enableAutoscaler {
		// Create a autoscaler
		if err := deployment.CreateAutoscaler(sp, autoscalerInfo); err != nil {
			if errors.IsAlreadyExists(err) {
				response.Conflict(req.Request, resp.ResponseWriter, err)
			} else if errors.IsConflict(err) {
				response.Conflict(req.Request, resp.ResponseWriter, err)
			} else if errors.IsInvalid(err) {
				response.BadRequest(req.Request, resp.ResponseWriter, err)
			} else {
				response.InternalServerError(req.Request, resp.ResponseWriter, err)
			}
			return

		}
	} else {
		// Delete a autoscaler
		if err := deployment.DeleteAutoscaler(sp, autoscalerInfo); err != nil {
			if errors.IsNotFound(err) {
				response.NotFound(req.Request, resp.ResponseWriter, err)
			} else if errors.IsForbidden(err) {
				response.NotAcceptable(req.Request, resp.ResponseWriter, err)
			} else {
				response.InternalServerError(req.Request, resp.ResponseWriter, err)
			}
			return
		}
	}
	// find the deployment id by deployment name
	deployment, err := backend.FindDeploymentByName(session, autoscalerInfo.ScaleTargetRefName)
	if err != nil {
		response.NotFound(req.Request, resp.ResponseWriter, fmt.Errorf("deployment not found"))
		return
	}

	// Update autoscalerInfo
	deployment.IsEnableAutoscale = enableAutoscaler
	deployment.AutoscalerInfo = autoscalerInfo
	modifier := bson.M{
		"$set": deployment,
	}
	query := bson.M{
		"_id": deployment.ID,
	}
	if err := session.Update(entity.DeploymentCollectionName, query, modifier); err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	resp.WriteHeaderAndEntity(http.StatusAccepted, deployment)
}
