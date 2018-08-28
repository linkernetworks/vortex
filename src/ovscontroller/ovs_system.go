package ovscontroller

import (
	"net"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/networkcontroller"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

func DumpPorts(sp *serviceprovider.Container, nodeName string, bridgeName string) ([]entity.OVSPortStat, error) {
	nodeIP, err := sp.KubeCtl.GetNodeInternalIP(nodeName)
	if err != nil {
		return nil, err
	}

	nodeAddr := net.JoinHostPort(nodeIP, networkcontroller.DEFAULT_CONTROLLER_PORT)
	nc, err := networkcontroller.New(nodeAddr)

	return nc.DumpOVSPorts(bridgeName)
}
