package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	PodCollectionName string = "pods"
)

type Container struct {
	Name    string   `bson:"name" json:"name"`
	Image   string   `bson:"image" json:"image"`
	Command []string `bson:"command" json:"command"`
}

type PodVolume struct {
	Name      string `bson:"name" json:"name"`
	MountPath string `bson:"mountPath" json:"mountPath"`
}

type Pod struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name       string        `bson:"name" json:"name"`
	Containers []Container   `bson:"containers" json:"containers"`
	CreatedAt  *time.Time    `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	Volumes    []PodVolume   `bson:"volumes,omitempty" json:"volumes"`
}

//GetCollection - get model mongo collection name.
func (m Pod) GetCollection() string {
	return PodCollectionName
}
