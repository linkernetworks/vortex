package backend

import (
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/entity"
	"gopkg.in/mgo.v2/bson"
)

// FindDeploymentByName will find deployment by name
func FindDeploymentByName(session *mongo.Session, Name string) (entity.Deployment, error) {
	var retDeployment entity.Deployment
	if err := session.FindOne(
		entity.DeploymentCollectionName,
		bson.M{"name": Name},
		&retDeployment,
	); err != nil {
		return entity.Deployment{}, err
	}
	return retDeployment, nil
}
