package server

import (
	"fmt"
	"math"
	"strconv"

	"github.com/linkernetworks/utils/timeutils"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/kubeutils"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/net/http/query"
	np "github.com/linkernetworks/vortex/src/networkprovider"
	"github.com/linkernetworks/vortex/src/web"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func createNetworkHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	network := entity.Network{}
	if err := req.ReadEntity(&network); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	// overwrite the bridge name
	network.BridgeName = np.GenerateBridgeName(string(network.Type), network.Name)

	if err := sp.Validator.Struct(network); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	session := sp.Mongo.NewSession()
	defer session.Close()
	session.C(entity.NetworkCollectionName).EnsureIndex(
		mgo.Index{
			Key:    []string{"name"},
			Unique: true,
		})

	networkProvider, err := np.GetNetworkProvider(&network)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := networkProvider.CreateNetwork(sp); err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
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

func listNetworkHandler(ctx *web.Context) {
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

	networks := []entity.Network{}
	var c = session.C(entity.NetworkCollectionName)
	var q *mgo.Query

	selector := bson.M{}
	q = c.Find(selector).Sort("_id").Skip((page - 1) * pageSize).Limit(pageSize)

	if err := q.All(&networks); err != nil {
		switch err {
		case mgo.ErrNotFound:
			response.NotFound(req.Request, resp.ResponseWriter, err)
			return
		default:
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
	}

	count, err := session.Count(entity.NetworkCollectionName, bson.M{})
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))
	resp.AddHeader("X-Total-Count", strconv.Itoa(count))
	resp.AddHeader("X-Total-Pages", strconv.Itoa(totalPages))
	resp.WriteEntity(networks)
}

func getNetworkHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	id := req.PathParameter("id")

	session := sp.Mongo.NewSession()
	defer session.Close()
	c := session.C(entity.NetworkCollectionName)

	var network entity.Network
	if err := c.FindId(bson.ObjectIdHex(id)).One(&network); err != nil {
		switch err {
		case mgo.ErrNotFound:
			response.NotFound(req.Request, resp.ResponseWriter, err)
			return
		default:
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
	}
	resp.WriteEntity(network)
}

func getNetworkStatusHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	id := req.PathParameter("id")

	session := sp.Mongo.NewSession()
	defer session.Close()
	c := session.C(entity.NetworkCollectionName)

	var network entity.Network
	if err := c.FindId(bson.ObjectIdHex(id)).One(&network); err != nil {
		switch err {
		case mgo.ErrNotFound:
			response.NotFound(req.Request, resp.ResponseWriter, err)
			return
		default:
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
	}

	//FIXME since the network doesn't have the namespace but Pod has..
	ret, err := kubeutils.GetNonCompletedPods(sp, bson.M{"networks.name": network.Name})
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	nameList := []string{}
	for _, v := range ret {
		nameList = append(nameList, v.Name)
	}
	resp.WriteEntity(nameList)
}

func deleteNetworkHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	id := req.PathParameter("id")

	session := sp.Mongo.NewSession()
	defer session.Close()
	c := session.C(entity.NetworkCollectionName)

	var network entity.Network
	if err := c.FindId(bson.ObjectIdHex(id)).One(&network); err != nil {
		switch err {
		case mgo.ErrNotFound:
			response.NotFound(req.Request, resp.ResponseWriter, err)
			return
		default:
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
	}

	networkProvider, err := np.GetNetworkProvider(&network)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := networkProvider.DeleteNetwork(sp); err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := session.Remove(entity.NetworkCollectionName, "_id", bson.ObjectIdHex(id)); err != nil {
		switch err {
		case mgo.ErrNotFound:
			response.NotFound(req.Request, resp.ResponseWriter, err)
			return
		default:
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
	}

	resp.WriteEntity(ActionResponse{
		Error:   false,
		Message: "Delete success",
	})
}
