package networkprovider

import (
	"fmt"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

type fakeNetworkProvider struct {
	entity.Network
}

func (fnp fakeNetworkProvider) CreateNetwork(sp *serviceprovider.Container) error {
	if !fnp.IsDPDKPort {
		return fmt.Errorf("fail to validate but don't worry, I'm fake network")
	}
	return nil
}

func (fnp fakeNetworkProvider) DeleteNetwork(sp *serviceprovider.Container) error {
	if !fnp.IsDPDKPort {
		return fmt.Errorf("fail to delete network but don't worry, I'm fake network")
	}
	return nil
}
