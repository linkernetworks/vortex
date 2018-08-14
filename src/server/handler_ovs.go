package server

import (
	"fmt"

	"github.com/linkernetworks/vortex/src/entity"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/net/http/query"
	"github.com/linkernetworks/vortex/src/ovscontroller"
	"github.com/linkernetworks/vortex/src/web"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getOVSPortStatsHandler(ctx *web.Context) {
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
	//Get the parameter
	query := query.New(req.Request.URL.Query())
	nodeName, exist := query.Str("nodeName")
	if !exist {
		response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("The nodeName must not be empty"))
		return
	}

	bridgeName, exist := query.Str("bridgeName")
	if !exist {
		response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("The bridgeName must not be empty"))
		return
	}

	fmt.Println(nodeName, bridgeName)
	portStats, err := ovscontroller.DumpPorts(sp, nodeName, bridgeName)
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	resp.WriteEntity(portStats)
}
