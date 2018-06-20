package entity

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	StorageProviderCollectionName string = "storage_provider"
)

type NFSStorageProvider struct {
	IP   string `bson:"ip" json:"ip"`
	PATH string `bson:"path" json:"path"`
}

type StorageProvider struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Type        string        `bson:"type" json:"type"`
	DisplayName string        `bson:"displayName" json:"displayName"`
	NFSStorageProvider
}

//GetCollection - get model mongo collection name.
func (m StorageProvider) GetCollection() string {
	return StorageProviderCollectionName
}
