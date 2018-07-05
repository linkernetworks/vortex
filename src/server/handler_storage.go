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
	"github.com/linkernetworks/vortex/src/storageprovider"
	"github.com/linkernetworks/vortex/src/web"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func createStorage(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	storage := entity.Storage{}
	if err := req.ReadEntity(&storage); err != nil {
		logger.Error(err)
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	storage.ID = bson.NewObjectId()
	storage.CreatedAt = timeutils.Now()

	session := sp.Mongo.NewSession()
	session.C(entity.StorageCollectionName).EnsureIndex(mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	})
	defer session.Close()
	// Check whether this name has been used
	storageProvider, err := storageprovider.GetStorageProvider(&storage)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := storageProvider.ValidateBeforeCreating(sp, &storage); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := storageProvider.CreateStorage(sp, &storage); err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := session.Insert(entity.StorageCollectionName, &storage); err != nil {
		if mgo.IsDup(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Storage Provider Name: %s already existed", storage.Name))
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

func listStorage(ctx *web.Context) {
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

	storageProviders := []entity.Storage{}

	var c = session.C(entity.StorageCollectionName)
	var q *mgo.Query

	selector := bson.M{}
	q = c.Find(selector).Sort("_id").Skip((page - 1) * pageSize).Limit(pageSize)

	if err := q.All(&storageProviders); err != nil {
		if err == mgo.ErrNotFound {
			response.NotFound(req.Request, resp.ResponseWriter, err)
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}

	count, err := session.Count(entity.StorageCollectionName, bson.M{})
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))
	resp.AddHeader("X-Total-Count", strconv.Itoa(count))
	resp.AddHeader("X-Total-Pages", strconv.Itoa(totalPages))
	resp.WriteEntity(storageProviders)
}

func deleteStorage(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	id := req.PathParameter("id")

	session := sp.Mongo.NewSession()
	defer session.Close()
	c := session.C(entity.StorageCollectionName)

	var storage entity.Storage
	if err := c.FindId(bson.ObjectIdHex(id)).One(&storage); err != nil {
		if err == mgo.ErrNotFound {
			response.NotFound(req.Request, resp.ResponseWriter, err)
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}

	storageProvider, err := storageprovider.GetStorageProvider(&storage)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := storageProvider.DeleteStorage(sp, &storage); err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := session.Remove(entity.StorageCollectionName, "_id", bson.ObjectIdHex(id)); err != nil {
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
