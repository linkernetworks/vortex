package networkprovider

import (
	"fmt"
	"net"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/networkcontroller"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"gopkg.in/mgo.v2/bson"
)

type KernelSpaceNetworkProvider struct {
	nodes []entity.Node
}

func (knp KernelSpaceNetworkProvider) ValidateBeforeCreating(sp *serviceprovider.Container, network *entity.Network) error {
	session := sp.Mongo.NewSession()
	defer session.Close()

	// Check whether VLAN Tag is 0~4095
	for _, tag := range network.VLANTags {
		if tag < 0 || tag > 4095 {
			return fmt.Errorf("The vlangTag %d should between 0 and 4095", tag)
		}
	}

	// TODO Check dpdk port
	for _, node := knp.nodes {
		node
	}

	q := bson.M{
		"networks.name": network.Name,
	}
	n, err := session.Count(entity.NetworkCollectionName, q)
	if n >= 1 {
		return fmt.Errorf("The network name: %s is exist.", network.Name)
	} else if err != nil {
		return err
	}
	return nil
}

func (knp KernelSpaceNetworkProvider) CreateNetwork(sp *serviceprovider.Container, network *entity.Network) error {
	for _, node := range knp.nodes {
		nodeIP, err := sp.KubeCtl.GetNodeExternalIP(node.Name)
		if err != nil {
			return err
		}
		if node.PhyInterface.IsDPDKPort {
			if err := createOVSDPDKNetwork(nodeIP, node.Name, node.PhyInterface, network.VLANTags); err != nil {
				return err
			}
		} else {
			if err := createOVSUserspaceNetwork(nodeIP, node.Name, node.PhyInterface, network.VLANTags); err != nil {
				return err
			}
		}
	}
	return nil
}

func (knp KernelSpaceNetworkProvider) DeleteNetwork(sp *serviceprovider.Container, network *entity.Network) error {
	for _, node := range knp.nodes {
		nodeIP, err := sp.KubeCtl.GetNodeExternalIP(node.Name)
		if err != nil {
			return err
		}
		if err := deleteOVSUserspaceNetwork(nodeIP, network.BridgeName); err != nil {
			return err
		}
	}
	return nil
}

func createOVSDPDKNetwork(nodeIP string, bridgeName string, phyIface entity.PhyInterface, vlanTags []int32) error {
	nodeAddr := net.JoinHostPort(nodeIP, networkcontroller.DEFAULT_CONTROLLER_PORT)
	nc, err := networkcontroller.New(nodeAddr)
	if err != nil {
		return err
	}
	return nc.CreateOVSDPDKNetwork(bridgeName, phyIface, vlanTags)
}

func createOVSUserspaceNetwork(nodeIP string, bridgeName string, phyIface entity.PhyInterface, vlanTags []int32) error {
	nodeAddr := net.JoinHostPort(nodeIP, networkcontroller.DEFAULT_CONTROLLER_PORT)
	nc, err := networkcontroller.New(nodeAddr)
	if err != nil {
		return err
	}
	return nc.CreateOVSNetwork("netdev", bridgeName, phyIface, vlanTags)
}

func deleteOVSUserspaceNetwork(nodeIP string, bridgeName string) error {
	nodeAddr := net.JoinHostPort(nodeIP, networkcontroller.DEFAULT_CONTROLLER_PORT)
	nc, err := networkcontroller.New(nodeAddr)
	if err != nil {
		return err
	}
	return nc.DeleteOVSNetwork(bridgeName)
}
