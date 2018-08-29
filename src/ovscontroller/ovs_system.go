package ovscontroller

import (
	"net"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/networkcontroller"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

func DumpPorts(sp *serviceprovider.Container, nodeName string, bridgeName string) ([]entity.OVSPortInfo, error) {
	nodeIP, err := sp.KubeCtl.GetNodeInternalIP(nodeName)
	if err != nil {
		return nil, err
	}

	nodeAddr := net.JoinHostPort(nodeIP, networkcontroller.DEFAULT_CONTROLLER_PORT)
	nc, err := networkcontroller.New(nodeAddr)

	retPorts, err := nc.DumpOVSPorts(bridgeName)
	if err != nil {
		return nil, err
	}

	ports := []entity.OVSPortInfo{}
	for _, v := range retPorts {
		ports = append(ports, entity.OVSPortInfo{
			PortID:     v.ID,
			Name:       v.Name,
			MacAddress: v.MacAddr,
			Received: entity.OVSPortStats{
				Packets: v.Received.Packets,
				Bytes:   v.Received.Byte,
				Dropped: v.Received.Dropped,
				Errors:  v.Received.Errors,
			},
			Transmitted: entity.OVSPortStats{
				Packets: v.Transmitted.Packets,
				Bytes:   v.Transmitted.Byte,
				Dropped: v.Transmitted.Dropped,
				Errors:  v.Transmitted.Errors,
			},
		})
	}
	return ports, nil
}
