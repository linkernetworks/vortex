package server

import (
	"fmt"

	"github.com/linkernetworks/utils/timeutils"
	"github.com/linkernetworks/vortex/src/entity"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/web"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func createVolume(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	volume := entity.Volume{}
	if err := req.ReadEntity(&volume); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	session := sp.Mongo.NewSession()
	session.C(entity.VolumeCollectionName).EnsureIndex(mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	})
	defer session.Close()

	// Check the storageClass is existed
	couunt, err := session.Count(entity.StorageProviderCollectionName, bson.M{"displayName": volume.StorageProviderName})
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	} else if couunt < 1 {
		response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("The reference storage provider %s doesn't exist", volume.StorageProviderName))
		return
	}

	// Check whether this displayname has been used
	volume.ID = bson.NewObjectId()
	volume.CreatedAt = timeutils.Now()
	//Generate the metaName for furture use
	volume.MetaName = volume.GenerateMetaName()
	if err := session.Insert(entity.VolumeCollectionName, &volume); err != nil {
		if mgo.IsDup(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Storage Provider Name: %s already existed", volume.Name))
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