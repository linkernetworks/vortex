package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	NetworkCollectionName string = "networks"
)

type Network struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Type        string        `bson:"type" json:"type"`
	Name        string        `bson:"name" json:"name"`
	Clusterwise bool          `bson:"clusterwise" json:"clusterwise"`
	NodeName    string        `bson:"nodeName,omitempty" json:"nodeName,omitempty"`
	CreatedAt   *time.Time    `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	OVS         OVSNetwork    `bson:"ovs,omitempty" json:"ovs"`
	Fake        FakeNetwork   `json:"fake"` //FakeNetwork, for restful testing.
}

//GetCollection - get model mongo collection name.
func (m Network) GetCollection() string {
	return NetworkCollectionName
}
