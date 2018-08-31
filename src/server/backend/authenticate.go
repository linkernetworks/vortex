package backend

import (
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/utils"
	"gopkg.in/mgo.v2/bson"
)

// SecreKey is not a real secret key. Using the linker option to replace
var (
	SecretKey = "linkernetworks"
)

// Authenticate is a user authenticate function
func Authenticate(session *mongo.Session, credential entity.LoginCredential) (entity.User, bool, error) {
	authenticatedUser := entity.User{}
	if err := session.FindOne(
		entity.UserCollectionName,
		bson.M{"loginCredential.username": credential.Username},
		&authenticatedUser,
	); err != nil {
		return entity.User{}, false, err
	}
	hashedPassword := authenticatedUser.LoginCredential.Password
	if utils.CheckPasswordHash(credential.Password, hashedPassword) {
		return authenticatedUser, true, nil
	}
	return entity.User{}, false, nil
}
