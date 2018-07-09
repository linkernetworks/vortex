package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type NetworkType string

const (
	OVSKernelspaceNetworkType NetworkType = "system"
	OVSUserspaceNetworkType   NetworkType = "netdev"
	FakeNetworkType           NetworkType = "fake"
)

const (
	NetworkCollectionName string = "networks"
)

type PhyInterfaces struct {
	Name  string `bson:"name" json:"name"`
	PCIID string `bson:"pciID" json:"pciID"`
}

type Node struct {
	Name          string          `bson:"name" json:"name"`
	PhyInterfaces []PhyInterfaces `bson:"physicalInterface" json:"physicalInterface"`

	// Fake fields for restful testing
	FakeParameter string `json:"fakeParameter"`
	ShouldFail    bool   `json:"shoulFail"`
}

type Network struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Type       NetworkType   `bson:"type" json:"type"`
	IsDPDKPort bool          `bson:"isDPDKPort" json:"isDPDKPort"`
	Name       string        `bson:"name" json:"name"`
	VLANTags   []int32       `bson:"VLANTags" json:"VLANTags"`
	BridgeName string        `bson:"bridgeName" json:"bridgeName"`
	Nodes      []Node        `bson:"nodes" json:"nodes"`
	CreatedAt  *time.Time    `bson:"createdAt,omitempty" json:"createdAt,omitempty"`

	// Clusterwise  bool                `bson:"clusterwise" json:"clusterwise"`
	// NodeName     string              `bson:"nodeName,omitempty" json:"nodeName,omitempty"`
	// OVS          OVSNetwork          `bson:"ovs,omitempty" json:"ovs"`
	// OVSUserspace OVSUserspaceNetwork `bson:"ovsUserspace,omitempty" json:"ovsUserspace"`
	// Fake      FakeNetwork `json:"fake"` //FakeNetwork, for restful testing.
}

// GetCollection - get model mongo collection name.
func (m Network) GetCollection() string {
	return NetworkCollectionName
}
