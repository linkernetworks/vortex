package server

import (
	"fmt"
	"math"
	"strconv"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/utils/timeutils"
	"github.com/linkernetworks/vortex/src/entity"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/net/http/query"
	"github.com/linkernetworks/vortex/src/networkcontroller"
	"github.com/linkernetworks/vortex/src/web"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func createNetworkHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	network := entity.Network{}
	if err := req.ReadEntity(&network); err != nil {
		logger.Error(err)
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	session := sp.Mongo.NewSession()
	defer session.Close()
	session.C(entity.NetworkCollectionName).EnsureIndex(mgo.Index{
		Key:    []string{"bridgeName", "nodeName"},
		Unique: true,
	})

	// Check whether vlangTag is 0~4095
	for _, pp := range network.PhysicalPorts {
		for _, vlangTag := range pp.VlanTags {
			if vlangTag < 0 || vlangTag > 4095 {
				response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("the vlangTag %v in PhysicalPort %v should between 0 and 4095", pp.Name, vlangTag))
				return
			}
		}
	}

<<<<<<< HEAD
	nc, err := networkcontroller.New(sp.KubeCtl, network)
=======
	nc, err := networkcontroller.New(as.KubeCtl, network)
>>>>>>> Update networkcontrol for kubectl
	if err != nil {
		logger.Errorf("Failed to new network controller: %s", err.Error())
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := nc.CreateNetwork(); err != nil {
		logger.Errorf("Failed to create network: %s", err.Error())
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	network.ID = bson.NewObjectId()
	network.CreatedAt = timeutils.Now()
	if err := session.Insert(entity.NetworkCollectionName, &network); err != nil {
		if mgo.IsDup(err) {
			response.Conflict(req.Request, resp, fmt.Errorf("Network Name: %s already existed", network.BridgeName))
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

func listNetworkHandler(ctx *web.Context) {
	as, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

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

func getNetworkHandler(ctx *web.Context) {
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

func deleteNetworkHandler(ctx *web.Context) {
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
