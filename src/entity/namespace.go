package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// the const for NamespaceCollectionName
const (
	NamespaceCollectionName string = "namespaces"
)

// Namespace is the structure for namespace info
type Namespace struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id" validate:"-"`
	OwnerID   bson.ObjectId `bson:"ownerID,omitempty" json:"ownerID" validate:"-"`
	Name      string        `bson:"name" json:"name" validate:"required,k8sname"`
	CreatedAt *time.Time    `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"-"`
	CreatedBy User          `json:"createdBy" validate:"-"`
}

// GetCollection - get model mongo collection name.
func (m Namespace) GetCollection() string {
	return NamespaceCollectionName
}
