package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	NetworkCollectionName string = "networks"
)

type Network struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Type      string        `bson:"type" json:"type"`
	Name      string        `bson:"name" json:"name"`
	NodeName  string        `bson:"nodeName" json:"nodeName"`
	OVS       OVSNetwork    `bson:"ovs" json:"ovs"`
	CreatedAt *time.Time    `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
}

//GetCollection - get model mongo collection name.
func (m Network) GetCollection() string {
	return NetworkCollectionName
}
