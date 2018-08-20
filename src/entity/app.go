package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	// AppCollectionName is a const string for mongo collection name
	AppCollectionName string = "deployments"
)

// Deployment is the structure for deployment info
type App struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id" validate:"-"`
	Deployment Deployment    `bson:"deployment" json:"deployment" validate:"required"`
	Service    Service       `bson:"service" json:"service" validate:"required"`

	CreatedAt *time.Time `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"-"`
}

// GetCollection - get model mongo collection name.
func (m App) GetCollection() string {
	return AppCollectionName
}
