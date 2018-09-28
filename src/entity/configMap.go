package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// the const for ConfigMapCollectionName
const (
	ConfigMapCollectionName string = "configmap"
)

// ConfigMap is the structure for namespace info
type ConfigMap struct {
	ID        bson.ObjectId     `bson:"_id,omitempty" json:"id" validate:"-"`
	OwnerID   bson.ObjectId     `bson:"ownerID,omitempty" json:"ownerID" validate:"-"`
	Name      string            `bson:"name" json:"name" validate:"required,k8sname"`
	Namespace string            `bson:"namespace" json:"namespace" validate:"required"`
	Data      map[string]string `bson:"data,omitempty" json:"data,omitempty" validate:"-"`
	CreatedAt *time.Time        `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"-"`
	CreatedBy User              `json:"createdBy" validate:"-"`
}

// GetCollection - get model mongo collection name.
func (m ConfigMap) GetCollection() string {
	return ConfigMapCollectionName
}
