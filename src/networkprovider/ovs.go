package networkprovider

import (
	"fmt"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/networkcontroller"
	"github.com/linkernetworks/vortex/src/serviceprovider"

	"gopkg.in/mgo.v2/bson"
)

type OVSNetworkProvider struct {
	entity.OVSNetwork
}

func (ovs OVSNetworkProvider) ValidateBeforeCreating(sp *serviceprovider.Container, net entity.Network) error {
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

	q := bson.M{"nodeName": net.NodeName, "ovs.bridgeName": ovs.BridgeName}
	fmt.Println(q)
	//Check the bridge name, we can't have the same bridge name in the same node
	n, err := session.Count(entity.NetworkCollectionName, q)
	if n >= 1 {
		return fmt.Errorf("The bridge name %s is exist on the node %s\n, please use another bridge name", ovs.BridgeName, net.NodeName)
	} else if err != nil {
		return err
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
