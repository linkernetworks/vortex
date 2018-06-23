package server

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/utils/timeutils"
	"github.com/linkernetworks/vortex/src/entity"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/web"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func CreateStorageProvider(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	storageProvider := entity.StorageProvider{}
	if err := req.ReadEntity(&storageProvider); err != nil {
		logger.Error(err)
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	session := sp.Mongo.NewSession()
	session.C(entity.StorageProviderCollectionName).EnsureIndex(mgo.Index{
		Key:    []string{"displayName"},
		Unique: true,
	})
	defer session.Close()
	// Check whether this displayname has been used
	storageProvider.ID = bson.NewObjectId()
	storageProvider.CreatedAt = timeutils.Now()
	if err := session.Insert(entity.StorageProviderCollectionName, &storageProvider); err != nil {
		if mgo.IsDup(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Storage Provider Name: %s already existed", storageProvider.DisplayName))
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
