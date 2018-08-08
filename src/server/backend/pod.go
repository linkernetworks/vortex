package backend

import (
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/entity"
	"gopkg.in/mgo.v2/bson"
)

// FindPodByID will find pod by ID
func FindPodByID(session *mongo.Session, ID bson.ObjectId) (entity.Pod, error) {
	var pod entity.Pod
	if err := session.FindOne(
		entity.PodCollectionName,
		bson.M{"_id": ID},
		&pod,
	); err != nil {
		return entity.Pod{}, err
	}
	return pod, nil
}
