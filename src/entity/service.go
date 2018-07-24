package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// the const for ServiceCollectionName
const (
	ServiceCollectionName string = "services"
)

// ServicePort is the structure for service port
type ServicePort struct {
	Name       string `bson:"name" json:"name" validate:"required,k8sname"`
	Port       int32  `bson:"port" json:"port" validate:"required"`
	TargetPort int    `bson:"targetPort" json:"targetPort" validate:"required,max=65535,min=1"`
	NodePort   int32  `bson:"nodePort" json:"nodePort" validate:"max=32767,min=30000"`
}

// Service is the structure for service
type Service struct {
	ID        bson.ObjectId     `bson:"_id,omitempty" json:"id" validate:"-"`
	Name      string            `bson:"name" json:"name" validate:"required,k8sname"`
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
