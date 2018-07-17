package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	PodCollectionName string = "pods"
)

type Container struct {
	Name    string   `bson:"name" json:"name" validate:"required"`
	Image   string   `bson:"image" json:"image" validate:"required"`
	Command []string `bson:"command" json:"command" validate:"required,dive,required"`
}

type PodVolume struct {
	Name      string `bson:"name" json:"name" validate:"required"`
	MountPath string `bson:"mountPath" json:"mountPath" validate:"required"`
}

type Pod struct {
	ID         bson.ObjectId     `bson:"_id,omitempty" json:"id" validate:"-"`
	Name       string            `bson:"name" json:"name" validate:"required"`
	Namespace  string            `bson:"namespace" json:"namespace" validate:"required"`
	Labels     map[string]string `bson:"labels,omitempty" json:"labels" validate:"required,dive,keys,alphanum,endkeys,required,alphanum"`
	Containers []Container       `bson:"containers" json:"containers" validate:"required,dive,required"`
	Volumes    []PodVolume       `bson:"volumes,omitempty" json:"volumes" validate:"required,dive,required"`
	CreatedAt  *time.Time        `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"-"`
}

//GetCollection - get model mongo collection name.
func (m Pod) GetCollection() string {
	return PodCollectionName
}
