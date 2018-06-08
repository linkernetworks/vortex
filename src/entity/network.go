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
	DisplayName string        `bson:"displayName" json:"displayName"`
	BridgeName  string        `bson:"bridgeName" json:"bridgeName"`
	BridgeType  string        `bson:"bridgeType" json:"bridgeType"`
	Node        string        `bson:"node" json:"node"`
	Interface   string        `bson:"interface" json:"interface"`
	Ports       []int32       `bson:"ports" json:"ports"`
	MTU         int32         `bson:"maximumTransmissionUnit" MTC:"maximumTransmissionUnit"`
	CreatedAt   *time.Time    `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
}

//GetCollection - get model mongo collection name.
func (m Network) GetCollection() string {
	return NetworkCollectionName
}
