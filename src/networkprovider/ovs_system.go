package networkprovider

import (
	"net"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/networkcontroller"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

type kernelspaceNetworkProvider struct {
	networkName string
	bridgeName  string
	vlanTags    []int32
	nodes       []entity.Node
	isDPDKPort  bool
}

func (knp kernelspaceNetworkProvider) CreateNetwork(sp *serviceprovider.Container) error {
	if err := entity.ValidateVLANTags(knp.vlanTags); err != nil {
		return err
	}
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
