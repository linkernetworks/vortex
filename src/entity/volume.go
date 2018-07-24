package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
	corev1 "k8s.io/api/core/v1"
)

// The const for volume & PVC
const (
	VolumeCollectionName string = "volume"
	PVCNamePrefix        string = "pvc-"
)

// Volume is the structure. Users will create the Volume from the storage and
// they can use those volumes in their containers. In the kubernetes implementation, it's PVC
// So the Volume will create a PVC type and connect to a known StorageClass
type Volume struct {
	ID          bson.ObjectId                     `bson:"_id,omitempty" json:"id" validate:"-"`
	Name        string                            `bson:"name" json:"name" validate:"required"`
	StorageName string                            `bson:"storageName" json:"storageName" validate:"required"`
	AccessMode  corev1.PersistentVolumeAccessMode `bson:"accessMode" json:"accessMode" validate:"required"`
	Capacity    string                            `bson:"capacity" json:"capacity" validate:"required"`
	CreatedAt   *time.Time                        `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"-"`
}

//GetCollection - get model mongo collection name.
func (m Volume) GetCollection() string {
	return VolumeCollectionName
}

// GetPVCName will get pvc name
func (m Volume) GetPVCName() string {
	return PVCNamePrefix + m.ID.Hex()
}
