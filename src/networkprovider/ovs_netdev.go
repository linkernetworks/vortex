package networkprovider

import (
	"net"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/networkcontroller"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

type userspaceNetworkProvider struct {
	networkName string
	bridgeName  string
	vlanTags    []int32
	nodes       []entity.Node
	isDPDKPort  bool
}

func (unp userspaceNetworkProvider) CreateNetwork(sp *serviceprovider.Container) error {
	if err := entity.ValidateVLANTags(unp.vlanTags); err != nil {
		return err
	}
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
