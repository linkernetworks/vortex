package networkprovider

import (
	"fmt"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

type NetworkProvider interface {
	CreateNetwork(sp *serviceprovider.Container) error
	DeleteNetwork(sp *serviceprovider.Container) error
}

func GetNetworkProvider(network *entity.Network) (NetworkProvider, error) {
	switch network.Type {
	case entity.OVSKernelspaceNetworkType:
		return kernelspaceNetworkProvider{
			network.Name,
			network.BridgeName,
			network.VLANTags,
			network.Nodes,
			network.IsDPDKPort,
		}, nil
	case entity.OVSUserspaceNetworkType:
		return userspaceNetworkProvider{
			network.Name,
			network.BridgeName,
			network.VLANTags,
			network.Nodes,
			network.IsDPDKPort,
		}, nil
	case entity.FakeNetworkType:
		return fakeNetworkProvider{
			network.Name,
			network.BridgeName,
			network.VLANTags,
			network.Nodes,
			network.IsDPDKPort,
		}, nil
	default:
		return nil, fmt.Errorf("Unsupported Network Type %s", network.Type)
	}
}
