package server

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/linkernetworks/utils/timeutils"
	"github.com/linkernetworks/vortex/src/entity"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/net/http/query"
	"github.com/linkernetworks/vortex/src/server/backend"
	"github.com/linkernetworks/vortex/src/web"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func signUpUserHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	user := entity.User{}
	if err := req.ReadEntity(&user); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	encryptedPassword, err := backend.HashPassword(user.LoginCredential.Password)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}
	user.LoginCredential.Password = encryptedPassword

	user.LoginCredential.Username = strings.ToLower(user.LoginCredential.Username)

	// sign up user only can ba the role of user
	user.Role = "user"

	if err := sp.Validator.Struct(user); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	session := sp.Mongo.NewSession()
	// make username(email) to be a unique key
	session.C(entity.UserCollectionName).EnsureIndex(mgo.Index{
		Key:    []string{"loginCredential.username"},
		Unique: true,
	})
	defer session.Close()

	user.ID = bson.NewObjectId()
	user.CreatedAt = timeutils.Now()

	if err := session.Insert(entity.UserCollectionName, &user); err != nil {
		if mgo.IsDup(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Username: %s already existed", user.LoginCredential.Username))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}
	resp.WriteHeaderAndEntity(http.StatusCreated, user)
}

func verifyTokenHandler(ctx *web.Context) {
	_, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	userID, ok := req.Attribute("UserID").(string)
	if !ok {
		response.Unauthorized(req.Request, resp.ResponseWriter, fmt.Errorf("Unauthorized: User ID is empty"))
		return
	}
	http.Redirect(resp.ResponseWriter, req.Request, "/v1/users/"+userID, 303)
}

func signInUserHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	session := sp.Mongo.NewSession()
	defer session.Close()

	credential := entity.LoginCredential{}
	if err := req.ReadEntity(&credential); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := sp.Validator.Struct(credential); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	authenticatedUser, passed, err := backend.Authenticate(session, credential)
	if err != nil {
		switch err {
		case mgo.ErrNotFound:
			response.Unauthorized(req.Request, resp.ResponseWriter, fmt.Errorf("Unauthorized: Failed to login. Incorrect authentication credentials"))
			return
		default:
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
	}

	// when authenticating not pass
	if !passed {
		response.Unauthorized(req.Request, resp.ResponseWriter, fmt.Errorf("Unauthorized: Failed to login. Incorrect authentication credentials"))
		return
	}

	// Passed
	tokenString, err := backend.GenerateToken(authenticatedUser.ID.Hex(), authenticatedUser)
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	resp.WriteEntity(response.ActionResponse{
		Error:   false,
		Message: tokenString,
	})
}

func createUserHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	user := entity.User{}
	if err := req.ReadEntity(&user); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	encryptedPassword, err := backend.HashPassword(user.LoginCredential.Password)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}
	user.LoginCredential.Password = encryptedPassword

	user.LoginCredential.Username = strings.ToLower(user.LoginCredential.Username)

	if err := sp.Validator.Struct(user); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	session := sp.Mongo.NewSession()
	// make username(email) to be a unique key
	session.C(entity.UserCollectionName).EnsureIndex(mgo.Index{
		Key:    []string{"loginCredential.username"},
		Unique: true,
	})
	defer session.Close()

	user.ID = bson.NewObjectId()
	user.CreatedAt = timeutils.Now()

	if err := session.Insert(entity.UserCollectionName, &user); err != nil {
		if mgo.IsDup(err) {
			response.Conflict(req.Request, resp.ResponseWriter, fmt.Errorf("Username: %s already existed", user.LoginCredential.Username))
		} else {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
		}
		return
	}
	resp.WriteEntity(user)
}

func deleteUserHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	id := req.PathParameter("id")

	session := sp.Mongo.NewSession()
	defer session.Close()

	user := entity.User{}
	if err := session.FindOne(
		entity.UserCollectionName,
		bson.M{"_id": bson.ObjectIdHex(id)},
		&user,
	); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := session.Remove(entity.UserCollectionName, "_id", bson.ObjectIdHex(id)); err != nil {
		switch err {
		case mgo.ErrNotFound:
			response.NotFound(req.Request, resp.ResponseWriter, err)
			return
		default:
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
	}

	resp.WriteEntity(response.ActionResponse{
		Error:   false,
		Message: "User Deleted Success",
	})
}

// listUserHandler is to list all registered users on vortex server. The role must to have admin permission to access it.
func listUserHandler(ctx *web.Context) {
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

	users := []entity.User{}
	var c = session.C(entity.UserCollectionName)
	var q *mgo.Query

	selector := bson.M{}
	q = c.Find(selector).Sort("_id").Skip((page - 1) * pageSize).Limit(pageSize)

	if err := q.All(&users); err != nil {
		switch err {
		case mgo.ErrNotFound:
			response.NotFound(req.Request, resp.ResponseWriter, err)
			return
		default:
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
	}

	count, err := session.Count(entity.UserCollectionName, bson.M{})
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))
	resp.AddHeader("X-Total-Count", strconv.Itoa(count))
	resp.AddHeader("X-Total-Pages", strconv.Itoa(totalPages))
	resp.WriteEntity(users)
}

func getUserHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	id := req.PathParameter("id")

	session := sp.Mongo.NewSession()
	defer session.Close()
	c := session.C(entity.UserCollectionName)

	user := entity.User{}
	if err := c.FindId(bson.ObjectIdHex(id)).One(&user); err != nil {
		switch err {
		case mgo.ErrNotFound:
			response.NotFound(req.Request, resp.ResponseWriter, err)
			return
		default:
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
	}
	resp.WriteEntity(user)
}
