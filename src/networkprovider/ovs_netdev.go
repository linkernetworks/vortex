package networkprovider

import (
	"net"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/networkcontroller"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

type userspaceNetworkProvider struct {
	entity.Network
}

func (unp userspaceNetworkProvider) CreateNetwork(sp *serviceprovider.Container) error {
	if err := entity.ValidateVLANTags(unp.VLANTags); err != nil {
		return err
	}
	for _, node := range unp.Nodes {
		nodeIP, err := sp.KubeCtl.GetNodeExternalIP(node.Name)
		if err != nil {
			return err
		}
		if unp.IsDPDKPort {
			if err := createOVSDPDKNetwork(
				nodeIP,
				unp.BridgeName,
				node.PhyInterfaces,
				unp.VLANTags,
			); err != nil {
				return err
			}
		} else {
			if err := createOVSUserspaceNetwork(
				nodeIP,
				unp.BridgeName,
				node.PhyInterfaces,
				unp.VLANTags,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func (unp userspaceNetworkProvider) DeleteNetwork(sp *serviceprovider.Container) error {
	for _, node := range unp.Nodes {
		nodeIP, err := sp.KubeCtl.GetNodeExternalIP(node.Name)
		if err != nil {
			return err
		}
		if err := deleteOVSUserspaceNetwork(
			nodeIP,
			unp.BridgeName,
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
