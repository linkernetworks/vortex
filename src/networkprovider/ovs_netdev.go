package networkprovider

import (
	"fmt"
	"net"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/networkcontroller"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"gopkg.in/mgo.v2/bson"
)

type OVSUserspaceNetworkProvider struct {
	entity.OVSUserspaceNetwork
}

func (ovsu OVSUserspaceNetworkProvider) ValidateBeforeCreating(sp *serviceprovider.Container, network *entity.Network) error {
	session := sp.Mongo.NewSession()
	defer session.Close()
	// FIXME validate both dpdk or userspace datapath
	// Check whether vlangTag is 0~4095
	for _, pp := range ovsu.DPDKPhysicalPorts {
		for _, vlangTag := range pp.VlanTags {
			if vlangTag < 0 || vlangTag > 4095 {
				return fmt.Errorf("The vlangTag %v in PhysicalPort %v should between 0 and 4095", pp.Name, vlangTag)
			}
		}
	}

	q := bson.M{}
	if network.Clusterwise {
		//Only check the bridge name
		q = bson.M{"ovsUserspace.bridgeName": ovsu.BridgeName}
	} else {
		q = bson.M{"nodeName": network.NodeName, "ovsUserspace.bridgeName": ovsu.BridgeName}
	}

	n, err := session.Count(entity.NetworkCollectionName, q)
	if n >= 1 {
		return fmt.Errorf("The bridge name %s is exist, please check your cluster type and reassign another bridge name", ovsu.BridgeName)
	} else if err != nil {
		return err
	}
	return nil
}

func (ovsu OVSUserspaceNetworkProvider) CreateNetwork(sp *serviceprovider.Container, network *entity.Network) error {
	if network.Clusterwise {
		nodes, _ := sp.KubeCtl.GetNodes()
		for _, v := range nodes {
			nodeIP, err := sp.KubeCtl.GetNodeExternalIP(v.GetName())
			if err != nil {
				return err
			}
			// TODO if dpdk==true
			if err := createOVSDPDKNetwork(nodeIP, network.OVSUserspace.BridgeName, network.OVSUserspace.DPDKPhysicalPorts); err != nil {
				return err
			}
			// TODO else dpdk==false
			// if err := createOVSUserspaceNetwork(nodeIP, network.OVS.BridgeName, network.OVS.PhysicalPorts); err != nil {
			// 	return err
			// }
		}
		return nil
	}
	nodeIP, err := sp.KubeCtl.GetNodeExternalIP(network.NodeName)
	if err != nil {
		return err
	}
	// TODO if dpdk==true
	return createOVSDPDKNetwork(nodeIP, network.OVSUserspace.BridgeName, network.OVSUserspace.DPDKPhysicalPorts)
	// TODO else dpdk==false
	// return createOVSUserspaceNetwork(nodeIP, network.OVS.BridgeName, network.OVS.PhysicalPorts)
}

func (ovsu OVSUserspaceNetworkProvider) DeleteNetwork(sp *serviceprovider.Container, network *entity.Network) error {
	if network.Clusterwise {
		nodes, _ := sp.KubeCtl.GetNodes()
		for _, v := range nodes {
			nodeIP, err := sp.KubeCtl.GetNodeExternalIP(v.GetName())
			if err != nil {
				return err
			}
			if err := deleteOVSUserspaceNetwork(nodeIP, ovsu.BridgeName); err != nil {
				return err
			}
		}
		return nil
	}

	nodeIP, err := sp.KubeCtl.GetNodeExternalIP(network.NodeName)
	if err != nil {
		return err
	}
	return deleteOVSUserspaceNetwork(nodeIP, ovsu.BridgeName)
}

func createOVSDPDKNetwork(nodeIP string, bridgeName string, ports []entity.DPDKPhysicalPort) error {
	nodeAddr := net.JoinHostPort(nodeIP, networkcontroller.DEFAULT_CONTROLLER_PORT)
	nc, err := networkcontroller.New(nodeAddr)
	if err != nil {
		return err
	}
	return nc.CreateOVSDPDKNetwork(bridgeName, ports)
}

func createOVSUserspaceNetwork(nodeIP string, bridgeName string, ports []entity.PhysicalPort) error {
	nodeAddr := net.JoinHostPort(nodeIP, networkcontroller.DEFAULT_CONTROLLER_PORT)
	nc, err := networkcontroller.New(nodeAddr)
	if err != nil {
		return err
	}
	return nc.CreateOVSUserpsaceNetwork(bridgeName, ports)
}

func deleteOVSUserspaceNetwork(nodeIP string, bridgeName string) error {
	nodeAddr := net.JoinHostPort(nodeIP, networkcontroller.DEFAULT_CONTROLLER_PORT)
	nc, err := networkcontroller.New(nodeAddr)
	if err != nil {
		return err
	}
	return nc.DeleteOVSNetwork(bridgeName)
}
