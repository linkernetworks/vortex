package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type StorageType string

const (
	NFSStorageType  = "nfs"
	FakeStorageType = "fake"
)

const (
	StorageCollectionName string = "storage"
)

type Storage struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Type      StorageType   `bson:"type" json:"type"`
	Name      string        `bson:"name" json:"name"`
	CreatedAt *time.Time    `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	NFS       *NFSStorage   `bson:"nfs,omitempty" json:"nfs,omitempty"`
	Fake      *FakeStorage  `bson:"fake,omitempty" json:"fake,omitempty"` //FakeStorage, for restful testing.
}

//GetCollection - get model mongo collection name.
func (m Storage) GetCollection() string {
	return StorageCollectionName
}
