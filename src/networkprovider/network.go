package networkprovider

import (
	"fmt"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/linkernetworks/vortex/src/utils"
)

type NetworkProvider interface {
	CreateNetwork(sp *serviceprovider.Container) error
	DeleteNetwork(sp *serviceprovider.Container) error
}

func GetNetworkProvider(network *entity.Network) (NetworkProvider, error) {
	switch network.Type {
	case entity.OVSKernelspaceNetworkType:
		return kernelspaceNetworkProvider{
			*network,
		}, nil
	case entity.OVSUserspaceNetworkType:
		return userspaceNetworkProvider{
			*network,
		}, nil
	case entity.FakeNetworkType:
		return fakeNetworkProvider{
			*network,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported Network Type %s", network.Type)
	}
}

func GenerateBridgeName(datapathType, networkName string) string {
	tmp := fmt.Sprintf("%s%s", datapathType, networkName)
	str := utils.SHA256String(tmp)
	return fmt.Sprintf("ovs-%s-%s", datapathType, str[0:6])
}
