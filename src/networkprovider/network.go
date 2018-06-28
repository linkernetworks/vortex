package network

import (
	"fmt"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

type NetworkProvider interface {
	ValidateBeforeCreating(sp *serviceprovider.Container) error
}

func GetNetworkProvider(net *entity.Network) (NetworkProvider, error) {
	switch net.Type {
	case "ovs":
		return OVSNetworkProvider{net.OVS}, nil
	default:
		return nil, fmt.Errorf("Unsupported Network Type %s", net.Type)
	}
}
