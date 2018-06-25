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
	"github.com/linkernetworks/vortex/src/web"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func createStorageProvider(ctx *web.Context) {
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

func listStorageProvider(ctx *web.Context) {
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

	storageProviders := []entity.StorageProvider{}

	var c = session.C(entity.StorageProviderCollectionName)
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

	count, err := session.Count(entity.StorageProviderCollectionName, bson.M{})
	if err != nil {
		logger.Error(err)
	}
	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))
	resp.AddHeader("X-Total-Count", strconv.Itoa(count))
	resp.AddHeader("X-Total-Pages", strconv.Itoa(totalPages))
	resp.WriteEntity(storageProviders)
}
