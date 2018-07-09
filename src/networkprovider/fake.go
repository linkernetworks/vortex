package networkprovider

import (
	"fmt"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

type fakeNetworkProvider struct {
	networkName string
	bridgeName  string
	vlanTags    []int32
	nodes       []entity.Node
	isDPDKPort  bool
}

func (fnp fakeNetworkProvider) ValidateBeforeCreating(sp *serviceprovider.Container) error {
	for _, node := range fnp.nodes {
		if node.FakeParameter == "" {
			return fmt.Errorf("Fail to validate but don't worry, I'm fake network")
		}
	}
	return nil
}

func (fnp fakeNetworkProvider) CreateNetwork(sp *serviceprovider.Container) error {
	for _, node := range fnp.nodes {
		if node.ShouldFail {
			return fmt.Errorf("Fail to validate but don't worry, I'm fake network")
		}
	}
	return nil
}

func (fnp fakeNetworkProvider) DeleteNetwork(sp *serviceprovider.Container) error {
	for _, node := range fnp.nodes {
		if node.ShouldFail {
			return fmt.Errorf("Fail to delete network but don't worry, I'm fake network")
		}
	}
	return nil
}
