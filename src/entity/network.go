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
	BridgeName    string         `bson:"bridgeName" json:"bridgeName"`
	BridgeType    string         `bson:"bridgeType" json:"bridgeType"`
	NodeName      string         `bson:"nodeName" json:"nodeName"`
	PhysicalPorts []PhysicalPort `bson:"physicalPorts" json:"physicalPorts"`
	CreatedAt     *time.Time     `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
}

//GetCollection - get model mongo collection name.
func (m Network) GetCollection() string {
	return NetworkCollectionName
}
