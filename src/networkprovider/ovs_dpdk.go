package networkprovider

import (
	"fmt"
	"net"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/networkcontroller"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"gopkg.in/mgo.v2/bson"
)

type OVSDPDKNetworkProvider struct {
	entity.OVSDPDKNetwork
}

func (ovsdpdk OVSDPDKNetworkProvider) ValidateBeforeCreating(sp *serviceprovider.Container, network *entity.Network) error {
	session := sp.Mongo.NewSession()
	defer session.Close()
	//Check whether vlangTag is 0~4095
	for _, pp := range ovsdpdk.DPDKPhysicalPorts {
		for _, vlangTag := range pp.VlanTags {
			if vlangTag < 0 || vlangTag > 4095 {
				return fmt.Errorf("The vlangTag %v in PhysicalPort %v should between 0 and 4095", pp.Name, vlangTag)
			}
		}
	}

	q := bson.M{}
	if network.Clusterwise {
		//Only check the bridge name
		q = bson.M{"ovs.bridgeName": ovsdpdk.BridgeName}
	} else {
		q = bson.M{"nodeName": network.NodeName, "ovs.bridgeName": ovsdpdk.BridgeName}
	}

	n, err := session.Count(entity.NetworkCollectionName, q)
	if n >= 1 {
		return fmt.Errorf("The bridge name %s is exist, please check your cluster type and reassign another bridge name", ovsdpdk.BridgeName)
	} else if err != nil {
		return err
	}
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
