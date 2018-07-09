package networkprovider

import (
	"fmt"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

type FakeNetworkProvider struct {
	nodes []entity.Node
}

func (fnp FakeNetworkProvider) ValidateBeforeCreating(sp *serviceprovider.Container, net *entity.Network) error {
	for _, node := range fnp.nodes {
		if node.FakeParameter == "" {
			return fmt.Errorf("Fail to validate but don't worry, I'm fake network")
		}
	}
	return nil
}

func (fnp FakeNetworkProvider) CreateNetwork(sp *serviceprovider.Container, net *entity.Network) error {
	for _, node := range fnp.nodes {
		if node.ShouldFail {
			return fmt.Errorf("Fail to validate but don't worry, I'm fake network")
		}
	}
	return nil
}

func (fnp FakeNetworkProvider) DeleteNetwork(sp *serviceprovider.Container, net *entity.Network) error {
	for _, node := range fnp.nodes {
		if node.ShouldFail {
			return fmt.Errorf("Fail to delete network but don't worry, I'm fake network")
		}
	}
	return nil
}
