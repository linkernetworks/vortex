package networkprovider

import (
	"fmt"
	"net"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/networkcontroller"
	"github.com/linkernetworks/vortex/src/serviceprovider"

	"gopkg.in/mgo.v2/bson"
)

type OVSNetworkProvider struct {
	entity.OVSNetwork
}

func (ovs OVSNetworkProvider) ValidateBeforeCreating(sp *serviceprovider.Container, network entity.Network) error {
	session := sp.Mongo.NewSession()
	defer session.Close()
	//Check whether vlangTag is 0~4095
	for _, pp := range ovs.PhysicalPorts {
		for _, vlangTag := range pp.VlanTags {
			if vlangTag < 0 || vlangTag > 4095 {
				return fmt.Errorf("The vlangTag %v in PhysicalPort %v should between 0 and 4095", pp.Name, vlangTag)
			}
		}
	}

	q := bson.M{"nodeName": network.NodeName, "ovs.bridgeName": ovs.BridgeName}
	//Check the bridge name, we can't have the same bridge name in the same node
	n, err := session.Count(entity.NetworkCollectionName, q)
	if n >= 1 {
		return fmt.Errorf("The bridge name %s is exist on the node %s\n, please use another bridge name", ovs.BridgeName, network.NodeName)
	} else if err != nil {
		return err
	}
	return nil
}

func (ovs OVSNetworkProvider) CreateNetwork(sp *serviceprovider.Container, network entity.Network) error {
	nodeIP, err := sp.KubeCtl.GetNodeExternalIP(network.NodeName)
	if err != nil {
		return err
	}

	target := net.JoinHostPort(nodeIP, networkcontroller.DEFAULT_CONTROLLER_PORT)
	nc, err := networkcontroller.New(target)
	if err != nil {
		return err
	}

	return nc.CreateOVSNetwork(ovs.BridgeName, ovs.PhysicalPorts)
}
