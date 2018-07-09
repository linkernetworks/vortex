package networkprovider

import (
	"fmt"
	"net"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/networkcontroller"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"gopkg.in/mgo.v2/bson"
)

type kernelspaceNetworkProvider struct {
	networkName string
	bridgeName  string
	vlanTags    []int32
	nodes       []entity.Node
	isDPDKPort  bool
}

func (knp kernelspaceNetworkProvider) ValidateBeforeCreating(sp *serviceprovider.Container) error {
	session := sp.Mongo.NewSession()
	defer session.Close()

	// Check whether VLAN Tag is 0~4095
	for _, tag := range knp.vlanTags {
		if tag < 0 || tag > 4095 {
			return fmt.Errorf("The vlangTag %d should between 0 and 4095", tag)
		}
	}

	if knp.isDPDKPort != false {
		return fmt.Errorf("unsupport dpdk in kernel space datapath")
	}

	q := bson.M{
		"name": knp.networkName,
	}
	n, err := session.Count(entity.NetworkCollectionName, q)
	if n >= 1 {
		return fmt.Errorf("The network name: %s is exist.", knp.networkName)
	} else if err != nil {
		return err
	}
	return nil
}

func (knp kernelspaceNetworkProvider) CreateNetwork(sp *serviceprovider.Container) error {
	for _, node := range knp.nodes {
		nodeIP, err := sp.KubeCtl.GetNodeExternalIP(node.Name)
		if err != nil {
			return err
		}
		if err := createOVSNetwork(
			nodeIP,
			knp.bridgeName,
			node.PhyInterfaces,
			knp.vlanTags,
		); err != nil {
			return err
		}
	}
	return nil
}

func (knp kernelspaceNetworkProvider) DeleteNetwork(sp *serviceprovider.Container) error {
	for _, node := range knp.nodes {
		nodeIP, err := sp.KubeCtl.GetNodeExternalIP(node.Name)
		if err != nil {
			return err
		}
		if err := deleteOVSNetwork(
			nodeIP,
			knp.bridgeName,
		); err != nil {
			return err
		}
	}
	return nil
}

func createOVSNetwork(nodeIP string, bridgeName string, phyIfaces []entity.PhyInterface, vlanTags []int32) error {
	nodeAddr := net.JoinHostPort(nodeIP, networkcontroller.DEFAULT_CONTROLLER_PORT)
	nc, err := networkcontroller.New(nodeAddr)
	if err != nil {
		return err
	}
	return nc.CreateOVSNetwork("system", bridgeName, phyIfaces, vlanTags)
}

func deleteOVSNetwork(nodeIP string, bridgeName string) error {
	nodeAddr := net.JoinHostPort(nodeIP, networkcontroller.DEFAULT_CONTROLLER_PORT)
	nc, err := networkcontroller.New(nodeAddr)
	if err != nil {
		return err
	}
	return nc.DeleteOVSNetwork(bridgeName)
}
