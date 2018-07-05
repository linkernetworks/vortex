package networkprovider

import (
	"fmt"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

type NetworkProvider interface {
	ValidateBeforeCreating(sp *serviceprovider.Container, net *entity.Network) error
	CreateNetwork(sp *serviceprovider.Container, net *entity.Network) error
	DeleteNetwork(sp *serviceprovider.Container, net *entity.Network) error
}

func GetNetworkProvider(network *entity.Network) (NetworkProvider, error) {
	switch network.Type {
	case entity.OVSNetworkType:
		return OVSNetworkProvider{network.OVS}, nil
	case entity.FakeNetworkType:
		return FakeNetworkProvider{network.Fake}, nil
	default:
		return nil, fmt.Errorf("Unsupported Network Type %s", network.Type)
	}
}
