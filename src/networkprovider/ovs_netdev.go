package networkprovider

import (
	"fmt"
	"net"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/networkcontroller"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"gopkg.in/mgo.v2/bson"
)

type userspaceNetworkProvider struct {
	networkName string
	bridgeName  string
	vlanTags    []int32
	nodes       []entity.Node
	isDPDKPort  bool
}

func (unp userspaceNetworkProvider) ValidateBeforeCreating(sp *serviceprovider.Container) error {
	session := sp.Mongo.NewSession()
	defer session.Close()

	// Check whether VLAN Tag is 0~4095
	for _, tag := range unp.vlanTags {
		if tag < 0 || tag > 4095 {
			return fmt.Errorf("The vlangTag %d should between 0 and 4095", tag)
		}
	}

	q := bson.M{
		"name": unp.networkName,
	}
	n, err := session.Count(entity.NetworkCollectionName, q)
	if n >= 1 {
		return fmt.Errorf("The network name: %s is exist.", unp.networkName)
	} else if err != nil {
		return err
	}
	return nil
}

func (unp userspaceNetworkProvider) CreateNetwork(sp *serviceprovider.Container) error {
	for _, node := range unp.nodes {
		nodeIP, err := sp.KubeCtl.GetNodeExternalIP(node.Name)
		if err != nil {
			return err
		}
		if unp.isDPDKPort {
			if err := createOVSDPDKNetwork(
				nodeIP,
				unp.bridgeName,
				node.PhyInterfaces,
				unp.vlanTags,
			); err != nil {
				return err
			}
		} else {
			if err := createOVSUserspaceNetwork(
				nodeIP,
				unp.bridgeName,
				node.PhyInterfaces,
				unp.vlanTags,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func (unp userspaceNetworkProvider) DeleteNetwork(sp *serviceprovider.Container) error {
	for _, node := range unp.nodes {
		nodeIP, err := sp.KubeCtl.GetNodeExternalIP(node.Name)
		if err != nil {
			return err
		}
		if err := deleteOVSUserspaceNetwork(
			nodeIP,
			unp.bridgeName,
		); err != nil {
			return err
		}
	}
	return nil
}

func createOVSDPDKNetwork(nodeIP string, bridgeName string, phyIfaces []entity.PhyInterface, vlanTags []int32) error {
	nodeAddr := net.JoinHostPort(nodeIP, networkcontroller.DEFAULT_CONTROLLER_PORT)
	nc, err := networkcontroller.New(nodeAddr)
	if err != nil {
		return err
	}
	return nc.CreateOVSDPDKNetwork(bridgeName, phyIfaces, vlanTags)
}

func createOVSUserspaceNetwork(nodeIP string, bridgeName string, phyIfaces []entity.PhyInterface, vlanTags []int32) error {
	nodeAddr := net.JoinHostPort(nodeIP, networkcontroller.DEFAULT_CONTROLLER_PORT)
	nc, err := networkcontroller.New(nodeAddr)
	if err != nil {
		return err
	}
	return nc.CreateOVSNetwork("netdev", bridgeName, phyIfaces, vlanTags)
}

func deleteOVSUserspaceNetwork(nodeIP string, bridgeName string) error {
	nodeAddr := net.JoinHostPort(nodeIP, networkcontroller.DEFAULT_CONTROLLER_PORT)
	nc, err := networkcontroller.New(nodeAddr)
	if err != nil {
		return err
	}
	return nc.DeleteOVSNetwork(bridgeName)
}
