package backend

import (
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/entity"
	"gopkg.in/mgo.v2/bson"
)

func FindUserByUUID(session *mongo.Session, uuid string) (entity.User, error) {
	var user entity.User
	if err := session.FindOne(
		entity.UserCollectionName,
		bson.M{"uuid": uuid},
		&user,
	); err != nil {
		return entity.User{}, err
	}
	return user, nil
}
