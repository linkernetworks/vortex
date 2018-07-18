package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	ServiceCollectionName string = "services"
)

type ServicePort struct {
	Name       string `bson:"name" json:"name" validate:"required"`
	Port       int32  `bson:"port" json:"port" validate:"required"`
	TargetPort int    `bson:"targetPort" json:"targetPort" validate:"required,max=65535,min=1"`
	NodePort   int32  `bson:"nodePort" json:"nodePort" validate:"max=32767,min=30000"`
}

type Service struct {
	ID        bson.ObjectId     `bson:"_id,omitempty" json:"id" validate:"-"`
	Name      string            `bson:"name" json:"name" validate:"required"`
	Namespace string            `bson:"namespace" json:"namespace" validate:"required"`
	Type      string            `bson:"type" json:"type" validate:"oneof=ClusterIP NodePort"`
	Selector  map[string]string `bson:"selector" json:"selector" validate:"required"`
	Ports     []ServicePort     `bson:"ports" json:"ports" validate:"required"`
	CreatedAt *time.Time        `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"-"`
}

//GetCollection - get model mongo collection name.
func (m Service) GetCollection() string {
	return ServiceCollectionName
}
