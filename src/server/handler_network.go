package server

import (
	"fmt"

	"bitbucket.org/linkernetworks/aurora/src/web"
	"bitbucket.org/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/logger"
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
	defer session.Close()

	// Check whether this displayname has been used
	query := bson.M{"displayName": user.DisplayName}
	existed := entity.Network{}
	if err := session.FindOne(entity.NetworkCollectionName, query, &existed); err != nil {
		if err.Error() != mgo.ErrNotFound.Error() {
			logger.Error(err)
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
		if len(existed.ID) > 1 {
			response.Conflict(req.Request, resp, fmt.Errorf("displayName: %s already existed", user.DisplayName))
		}
	}

	// Check whether this Interface has been used
	query := bson.M{"node": user.Node, "interface": user.Interface}
	existed := entity.Network{}
	if err := session.FindOne(entity.NetworkCollectionName, query, &existed); err != nil {
		if err.Error() != mgo.ErrNotFound.Error() {
			logger.Error(err)
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
		if len(existed.ID) > 1 {
			response.Conflict(req.Request, resp, fmt.Errorf("interface %s on the Node has already be used", user.Interface, user.Node))
		}
	}

	// Check whether this bridge has been used
	query := bson.M{"bridge": user.Bridge}
	existed := entity.Network{}
	if err := session.FindOne(entity.NetworkCollectionName, query, &existed); err != nil {
		if err.Error() != mgo.ErrNotFound.Error() {
			logger.Error(err)
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
		if len(existed.ID) > 1 {
			response.Conflict(req.Request, resp, fmt.Errorf("bridge: %s already existed", user.Bridge))
		}
	}

	network.ID = bson.NewObjectId()
	network.CreatedAt = util.GetCurrentTimestamp()

	if err := session.Insert(entity.NetworkCollectionName, &user); err != nil {
		logger.Error(err)
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	resp.WriteEntity(ActionResponse{
		Error:   false,
		Message: "Create success",
	})
}
