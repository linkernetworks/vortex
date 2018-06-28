package network

import (
	"fmt"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/networkcontroller"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

type OVSNetworkProvider struct {
	entity.OVSNetwork
}

func (ovs OVSNetworkProvider) ValidateBeforeCreating(sp *serviceprovider.Container, net entity.Network) error {
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

func (ovs OVSNetworkProvider) CreateNetwork(sp *serviceprovider.Container, net entity.Network) error {
	nodeIP, err := sp.KubeCtl.GetNodeExternalIP(net.NodeName)
	if err != nil {
		return err
	}

	nc, err := networkcontroller.New(nodeIP + ":50051")
	if err != nil {
		return err
	}

	return nc.CreateOVSNetwork(ovs.BridgeName, ovs.PhysicalPorts)
}
