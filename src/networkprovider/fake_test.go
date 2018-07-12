package networkprovider

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linkernetworks/vortex/src/entity"
)

func TestFakeNetworkCreating(t *testing.T) {
	fake, err := GetNetworkProvider(&entity.Network{
		Type: entity.FakeNetworkType,
		Nodes: []entity.Node{
			entity.Node{
				FakeParameter: "yes",
				ShouldFail:    false,
			},
		},
	})
	assert.NoError(t, err)
	err = fake.CreateNetwork(nil)
	assert.NoError(t, err)
}

func TestFakeNetworkCreatingFail(t *testing.T) {
	fake, err := GetNetworkProvider(&entity.Network{
		Type: entity.FakeNetworkType,
		Nodes: []entity.Node{
			entity.Node{
				ShouldFail: true,
			},
		},
	})
	assert.NoError(t, err)
	err = fake.CreateNetwork(nil)
	assert.Error(t, err)
}

func TestFakeNetworkDelete(t *testing.T) {
	fake, err := GetNetworkProvider(&entity.Network{
		Type:  entity.FakeNetworkType,
		Nodes: []entity.Node{},
	})
	assert.NoError(t, err)
	err = fake.DeleteNetwork(nil)
	assert.NoError(t, err)
}

func TestFakeNetworkDeleteFail(t *testing.T) {
	fake, err := GetNetworkProvider(&entity.Network{
		Type: entity.FakeNetworkType,
		Nodes: []entity.Node{
			entity.Node{
				ShouldFail: true,
			},
		},
	})
	assert.NoError(t, err)
	err = fake.DeleteNetwork(nil)
	assert.Error(t, err)
}
