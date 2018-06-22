package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	NetworkCollectionName string = "networks"
)

type PhysicalPort struct {
	Name     string `bson:"name" json:"name"`
	MTU      int    `bson:"maximumTransmissionUnit" MTC:"maximumTransmissionUnit"`
	VlanTags []int  `bson:"vlanTag" MTC:"vlanTag"`
}

type Network struct {
	ID            bson.ObjectId  `bson:"_id,omitempty" json:"id"`
	Name          string         `bson:"name" json:"name"`
	BridgeType    string         `bson:"bridgeType" json:"bridgeType"`
	Node          string         `bson:"node" json:"node"`
	PhysicalPorts []PhysicalPort `bson:"physicalPorts" json:"physicalPorts"`
	CreatedAt     *time.Time     `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
}

//GetCollection - get model mongo collection name.
func (m Network) GetCollection() string {
	return NetworkCollectionName
}
