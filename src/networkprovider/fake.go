package networkprovider

import (
	"fmt"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

type FakeNetworkProvider struct {
	entity.FakeNetwork
}

func (fake FakeNetworkProvider) ValidateBeforeCreating(sp *serviceprovider.Container, net entity.Network) error {
	if fake.FakeParameter == "" {
		return fmt.Errorf("Fail to validate but don't worry, I'm fake network")
	}
	return nil
}

func (fake FakeNetworkProvider) CreateNetwork(sp *serviceprovider.Container, net entity.Network) error {
	if fake.IWantFail {
		return fmt.Errorf("Fail to create network but don't worry, I'm fake network")
	}
	return nil
}
