package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
	corev1 "k8s.io/api/core/v1"
)

const (
	VolumeCollectionName string = "volume"
)

/*
	Users will create the Volume from the storage and they can use those volumes in their containers
	In the kubernetes implementation, it's PVC
	So the Volume will create a PVC type and connect to a known StorageClass
*/
type Volume struct {
	ID          bson.ObjectId                     `bson:"_id,omitempty" json:"id"`
	Name        string                            `bson:"name" json:"name"`
	StorageName string                            `bson:"storageName" json:"storageName"`
	AccessMode  corev1.PersistentVolumeAccessMode `bson:"accessMode" json:"accessMode"`
	Capacity    string                            `bson:"capacity" json:"capacity"`
	MetaName    string                            `bson:"metaName" json:"metaName"` //For PVC metaname
	CreatedAt   *time.Time                        `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
}

//GetCollection - get model mongo collection name.
func (m Volume) GetCollection() string {
	return VolumeCollectionName
}

//GenerateMetaName - Generate a metaname for kubernetes PVC object
func (m Volume) GenerateMetaName() string {
	return "pvc-" + m.ID.Hex()
}
