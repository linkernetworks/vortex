package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// StorageType is string
type StorageType string

// The const for storage type
const (
	NFSStorageType  = "nfs"
	FakeStorageType = "fake"
)

// The const for StorageCollectionName
const (
	StorageCollectionName string = "storage"
)

// Storage is the Storage info
type Storage struct {
	ID               bson.ObjectId `bson:"_id,omitempty" json:"id" validate:"-"`
	Type             StorageType   `bson:"type" json:"type" validate:"required"`
	Name             string        `bson:"name" json:"name" validate:"required"`
	StorageClassName string        `bson:"storageClassName" json:"storageClassName" validate:"-"`
	IP               string        `bson:"ip" json:"ip" validate:"required,ipv4"`
	PATH             string        `bson:"path" json:"path" validate:"required"`
	Fake             *FakeStorage  `bson:"fake,omitempty" json:"fake,omitempty" validate:"-"` //FakeStorage, for restful testing.
	CreatedAt        *time.Time    `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"-"`
}

// GetCollection - get model mongo collection name.
func (m Storage) GetCollection() string {
	return StorageCollectionName
}
