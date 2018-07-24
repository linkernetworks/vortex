package entity

import (
	"time"

	"github.com/linkernetworks/vortex/src/errors"
	"gopkg.in/mgo.v2/bson"
)

// NetworkType is the string for network type
type NetworkType string

// These are const
const (
	OVSKernelspaceNetworkType NetworkType = "system"
	OVSUserspaceNetworkType   NetworkType = "netdev"
	FakeNetworkType           NetworkType = "fake"
)

// The const for NetworkCollectionName
const (
	NetworkCollectionName string = "networks"
)

// PhyInterface is the structure for physical interface
type PhyInterface struct {
	Name  string `bson:"name" json:"name" validate:"required"`
	PCIID string `bson:"pciID" json:"pciID" validate:"-"`
}

// Node is the structure for node info
type Node struct {
	Name          string         `bson:"name" json:"name" validate:"required"`
	PhyInterfaces []PhyInterface `bson:"physicalInterfaces" json:"physicalInterfaces" validate:"required,dive,required"`
}

// Network is the structure for Network info
type Network struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id" validate:"-"`
	Type       NetworkType   `bson:"type" json:"type" validate:"required"`
	IsDPDKPort bool          `bson:"isDPDKPort" json:"isDPDKPort" validate:"-"`
	Name       string        `bson:"name" json:"name" validate:"required"`
	VLANTags   []int32       `bson:"vlanTags" json:"vlanTags" validate:"required,dive,max=4095,min=0"`
	BridgeName string        `bson:"bridgeName" json:"bridgeName" validate:"-"`
	Nodes      []Node        `bson:"nodes" json:"nodes" validate:"required,dive,required"`
	CreatedAt  *time.Time    `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"-"`
}

// GetCollection - get model mongo collection name.
func (m Network) GetCollection() string {
	return NetworkCollectionName
}

// ValidateVLANTags will validate VLAN tags
func ValidateVLANTags(vlanTags []int32) error {
	for _, tag := range vlanTags {
		if tag < 0 || tag > 4095 {
			return errors.NewErrInvalidVLAN("VLAN tag should between 0 and 4095")
		}
	}
	return nil
}
