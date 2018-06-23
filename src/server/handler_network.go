package server

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/utils/timeutils"
	"github.com/linkernetworks/vortex/src/entity"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/net/http/query"
	"github.com/linkernetworks/vortex/src/web"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func CreateNetworkHandler(ctx *web.Context) {
	as, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	network := entity.Network{}

	if err := req.ReadEntity(&network); err != nil {
		logger.Error(err)
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	session := as.Mongo.NewSession()
	session.C(entity.NetworkCollectionName).EnsureIndex(mgo.Index{
		Key:    []string{"Name", "Node"},
		Unique: true,
	})
	defer session.Close()

	// Check whether vlangTag is 0~4095
	for _, pp := range network.PhysicalPorts {
		for _, vlangTag := range pp.VlanTags {
			if vlangTag < 0 || vlangTag > 4095 {
				response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("the vlangTag %v in PhysicalPort %v should between 0 and 4095", pp.Name, vlangTag))
				return
			}
		}
	}

	network.ID = bson.NewObjectId()
	network.CreatedAt = timeutils.Now()
	if err := session.Insert(entity.NetworkCollectionName, &network); err != nil {
		if mgo.IsDup(err) {
			response.Conflict(req.Request, resp, fmt.Errorf("Network Name: %s already existed", network.Name))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}

	resp.WriteEntity(ActionResponse{
		Error:   false,
		Message: "Create success",
	})
}

func ListNetworkHandler(ctx *web.Context) {
	as, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	var pageSize = 10
	query := query.New(req.Request.URL.Query())

	page, err := query.Int("page", 1)
	if err != nil {
		logger.Error(err)
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}
	pageSize, err = query.Int("page_size", pageSize)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	session := as.Mongo.NewSession()
	defer session.Close()

	networks := []entity.Network{}

	var c = session.C(entity.NetworkCollectionName)
	var q *mgo.Query

	selector := bson.M{}
	q = c.Find(selector).Sort("_id").Skip((page - 1) * pageSize).Limit(pageSize)

	if err := q.All(&networks); err != nil {
		if err == mgo.ErrNotFound {
			response.NotFound(req.Request, resp.ResponseWriter, err)
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}

	count, err := session.Count(entity.NetworkCollectionName, bson.M{})
	if err != nil {
		logger.Error(err)
	}
	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))
	resp.AddHeader("X-Total-Count", strconv.Itoa(count))
	resp.AddHeader("X-Total-Pages", strconv.Itoa(totalPages))
	resp.WriteEntity(networks)
}

func GetNetworkHandler(ctx *web.Context) {
	as, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	id := req.PathParameter("id")

	session := as.Mongo.NewSession()
	defer session.Close()
	c := session.C(entity.NetworkCollectionName)

	var network entity.Network
	if err := c.FindId(bson.ObjectIdHex(id)).One(&network); err != nil {
		if err == mgo.ErrNotFound {
			response.NotFound(req.Request, resp.ResponseWriter, err)
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}
	resp.WriteEntity(network)
}

func DeleteNetworkHandler(ctx *web.Context) {
	as, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	id := req.PathParameter("id")

	session := as.Mongo.NewSession()
	defer session.Close()

	if err := session.Remove(entity.NetworkCollectionName, "_id", bson.ObjectIdHex(id)); err != nil {
		if mgo.ErrNotFound == err {
			response.NotFound(req.Request, resp.ResponseWriter, err)
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}

	resp.WriteEntity(ActionResponse{
		Error:   false,
		Message: "Delete success",
	})
}

func UpdateNetworkHandler(ctx *web.Context) {
	as, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	id := req.PathParameter("id")

	session := as.Mongo.NewSession()
	defer session.Close()

	network := entity.Network{}
	q := bson.M{"_id": bson.ObjectIdHex(id)}
	if err := session.FindOne(entity.NetworkCollectionName, q, &network); err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			response.NotFound(req.Request, resp.ResponseWriter, fmt.Errorf("the network: %v doesn't exist", id))
			return
		}
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	updatedNetwork := entity.Network{}
	if err := req.ReadEntity(&updatedNetwork); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	checkNetwork := entity.Network{}
	checkNetwork.Name = updatedNetwork.Name
	if !reflect.DeepEqual(updatedNetwork, checkNetwork) {
		response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("only Network Namea can be changed"))
		return
	}

	if err := session.UpdateById(entity.NetworkCollectionName, network.ID, updatedNetwork); err != nil {
		if err == mgo.ErrNotFound {
			response.NotFound(req.Request, resp.ResponseWriter, err)
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}
	resp.WriteEntity(updatedNetwork)
}
