package networkprovider

import (
	"net"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/networkcontroller"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

type OVSDPDKNetworkProvider struct {
	entity.OVSDPDKNetwork
}

func (ovsdpdk OVSDPDKNetworkProvider) ValidateBeforeCreating(sp *serviceprovider.Container, network *entity.Network) error {
	return nil
}

func (ovsdpdk OVSDPDKNetworkProvider) CreateNetwork(sp *serviceprovider.Container, network *entity.Network) error {
	if network.Clusterwise {
		nodes, _ := sp.KubeCtl.GetNodes()
		for _, v := range nodes {
			nodeIP, err := sp.KubeCtl.GetNodeExternalIP(v.GetName())
			if err != nil {
				return err
			}
			if err := createOVSDPDKNetwork(nodeIP, network.OVS.BridgeName, network.OVSDPDK.DPDKPhysicalPorts); err != nil {
				return err
			}
		}
		return nil
	}
	nodeIP, err := sp.KubeCtl.GetNodeExternalIP(network.NodeName)
	if err != nil {
		return err
	}
	return createOVSDPDKNetwork(nodeIP, network.OVS.BridgeName, network.OVSDPDK.DPDKPhysicalPorts)
}

func (ovsdpdk OVSDPDKNetworkProvider) DeleteNetwork(sp *serviceprovider.Container, network *entity.Network) error {
	if network.Clusterwise {
		nodes, _ := sp.KubeCtl.GetNodes()
		for _, v := range nodes {
			nodeIP, err := sp.KubeCtl.GetNodeExternalIP(v.GetName())
			if err != nil {
				return err
			}
			if err := deleteOVSNetwork(nodeIP, ovsdpdk.BridgeName); err != nil {
				return err
			}
		}
		return nil
	}

	nodeIP, err := sp.KubeCtl.GetNodeExternalIP(network.NodeName)
	if err != nil {
		return err
	}
	return deleteOVSNetwork(nodeIP, ovsdpdk.BridgeName)
}

func createOVSDPDKNetwork(nodeIP string, bridgeName string, ports []entity.DPDKPhysicalPort) error {
	nodeAddr := net.JoinHostPort(nodeIP, networkcontroller.DEFAULT_CONTROLLER_PORT)
	nc, err := networkcontroller.New(nodeAddr)
	if err != nil {
		return err
	}
	return nc.CreateOVSDPDKNetwork(bridgeName, ports)
}

func deleteOVSDPDKNetwork(nodeIP string, bridgeName string) error {
	nodeAddr := net.JoinHostPort(nodeIP, networkcontroller.DEFAULT_CONTROLLER_PORT)
	nc, err := networkcontroller.New(nodeAddr)
	if err != nil {
		return err
	}
	return nc.DeleteOVSDPDKNetwork(bridgeName)
}
