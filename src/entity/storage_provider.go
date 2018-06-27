package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	StorageCollectionName string = "storage"
)

type NFSStorageSetting struct {
	IP   string `bson:"ip" json:"ip"`
	PATH string `bson:"path" json:"path"`
}

type Storage struct {
	ID                bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Type              string        `bson:"type" json:"type"`
	DisplayName       string        `bson:"displayName" json:"displayName"`
	CreatedAt         *time.Time    `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	NFSStorageSetting `bson:"nfs" json:"nfs"`
}

//GetCollection - get model mongo collection name.
func (m Storage) GetCollection() string {
	return StorageCollectionName
}
