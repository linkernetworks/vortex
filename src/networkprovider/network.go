package network

import (
	"fmt"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

type NetworkProvider interface {
	ValidateBeforeCreating(sp *serviceprovider.Container, net entity.Network) error
	CreateNetwork(sp *serviceprovider.Container, net entity.Network) error
}

func GetNetworkProvider(net *entity.Network) (NetworkProvider, error) {
	switch net.Type {
	case "ovs":
		return OVSNetworkProvider{net.OVS}, nil
	case "fake":
		return FakeNetworkProvider{net.Fake}, nil
	default:
		return nil, fmt.Errorf("Unsupported Network Type %s", net.Type)
	}
}
