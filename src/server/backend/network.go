package backend

import (
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/entity"
	"gopkg.in/mgo.v2/bson"
)

func FindNetworkByID(session *mongo.Session, ID bson.ObjectId) (entity.Network, error) {
	var network entity.Network
	if err := session.FindOne(
		entity.NetworkCollectionName,
		bson.M{"_id": ID},
		&network,
	); err != nil {
		return entity.Network{}, err
	}
	return network, nil
}
