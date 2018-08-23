package backend

import (
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/entity"
	"gopkg.in/mgo.v2/bson"
)

func FindUserByID(session *mongo.Session, ID bson.ObjectId) (entity.User, error) {
	var user entity.User
	if err := session.FindOne(
		entity.UserCollectionName,
		bson.M{"_id": ID},
		&user,
	); err != nil {
		return entity.User{}, err
	}
	return user, nil
}
