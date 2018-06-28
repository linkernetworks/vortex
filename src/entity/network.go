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
	CreatedAt *time.Time    `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	OVS       OVSNetwork    `bson:"ovs,omitempty" json:"ovs"`
	Fake      FakeNetwork   `bson:"fake, omitempty" json:"fake"` //FakeNetwork, for restful testing.
}

//GetCollection - get model mongo collection name.
func (m Network) GetCollection() string {
	return NetworkCollectionName
}
