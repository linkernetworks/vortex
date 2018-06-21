package server

import (
	"fmt"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/utils/timeutils"
	"github.com/linkernetworks/vortex/src/entity"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/web"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func CreateStorageProvider(ctx *web.Context) {
	as, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	storageProvider := entity.StorageProvider{}

	if err := req.ReadEntity(&storageProvider); err != nil {
		logger.Error(err)
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	session := as.Mongo.NewSession()
	defer session.Close()

	// Check whether this displayname has been used
	query := bson.M{"displayName": storageProvider.DisplayName}
	count, err := session.Count(entity.NetworkCollectionName, query)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logger.Error(err)
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	if count > 0 {
		response.Conflict(req.Request, resp, fmt.Errorf("displayName: %s already existed", storageProvider.DisplayName))
		return
	}

	storageProvider.ID = bson.NewObjectId()
	storageProvider.CreatedAt = timeutils.Now()
	if err := session.Insert(entity.StorageProviderCollectionName, &storageProvider); err != nil {
		logger.Error(err)
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	resp.WriteEntity(ActionResponse{
		Error:   false,
		Message: "Create success",
	})
}
