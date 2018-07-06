package server

import (
	"fmt"
	"math"
	"strconv"

	"github.com/linkernetworks/utils/timeutils"
	"github.com/linkernetworks/vortex/src/entity"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/net/http/query"
	"github.com/linkernetworks/vortex/src/volume"
	"github.com/linkernetworks/vortex/src/web"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func createVolume(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	v := entity.Volume{}
	if err := req.ReadEntity(&v); err != nil {
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
	couunt, err := session.Count(entity.StorageCollectionName, bson.M{"name": v.StorageName})
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	} else if couunt < 1 {
		response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("The reference storage provider %s doesn't exist", v.StorageName))
		return
	}

	// Check whether this name has been used
	v.ID = bson.NewObjectId()
	v.CreatedAt = timeutils.Now()
	v.MetaName = v.GenerateMetaName()
	//Generate the metaName for PVC meta name and we will use it future
	if err := volume.CreateVolume(sp, &v); err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	if err := session.Insert(entity.VolumeCollectionName, &v); err != nil {
		if mgo.IsDup(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Storage Provider Name: %s already existed", v.Name))
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

func deleteVolume(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	id := req.PathParameter("id")

	session := sp.Mongo.NewSession()
	defer session.Close()

	v := entity.Volume{}
	if err := session.FindOne(entity.VolumeCollectionName, bson.M{"_id": bson.ObjectIdHex(id)}, &v); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := volume.DeleteVolume(sp, &v); err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		fmt.Println("QQ-", err)
		return
	}

	if err := session.Remove(entity.VolumeCollectionName, "_id", bson.ObjectIdHex(id)); err != nil {
		if mgo.ErrNotFound == err {
			response.BadRequest(req.Request, resp.ResponseWriter, err)
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

func listVolume(ctx *web.Context) {
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

	volumes := []entity.Volume{}
	var c = session.C(entity.VolumeCollectionName)
	var q *mgo.Query

	selector := bson.M{}
	q = c.Find(selector).Sort("_id").Skip((page - 1) * pageSize).Limit(pageSize)

	if err := q.All(&volumes); err != nil {
		if err == mgo.ErrNotFound {
			response.NotFound(req.Request, resp.ResponseWriter, err)
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}

	count, err := session.Count(entity.VolumeCollectionName, bson.M{})
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))
	resp.AddHeader("X-Total-Count", strconv.Itoa(count))
	resp.AddHeader("X-Total-Pages", strconv.Itoa(totalPages))
	resp.WriteEntity(volumes)
}
