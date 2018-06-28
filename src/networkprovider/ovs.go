package network

import (
	"fmt"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

type OVSNetworkProvider struct {
	entity.OVSNetwork
}

func (ovs OVSNetworkProvider) ValidateBeforeCreating(sp *serviceprovider.Container) error {
	session := sp.Mongo.NewSession()
	defer session.Close()
	// Check whether vlangTag is 0~4095
	for _, pp := range ovs.PhysicalPorts {
		for _, vlangTag := range pp.VlanTags {
			if vlangTag < 0 || vlangTag > 4095 {
				return fmt.Errorf("the vlangTag %v in PhysicalPort %v should between 0 and 4095", pp.Name, vlangTag)
			}
		}
	}
	return nil
}
