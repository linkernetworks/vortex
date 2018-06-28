package networkprovider

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linkernetworks/vortex/src/entity"
)

func TestFakeNetworkValidateBeforeCreating(t *testing.T) {
	fake, err := GetNetworkProvider(&entity.Network{
		Type: "fake",
		Fake: entity.FakeNetwork{
			FakeParameter: "yes",
		},
	})
	assert.NoError(t, err)
	err = fake.ValidateBeforeCreating(nil, entity.Network{})
	assert.NoError(t, err)
}

func TestFakeNetworkCreating(t *testing.T) {
	fake, err := GetNetworkProvider(&entity.Network{
		Type: "fake",
		Fake: entity.FakeNetwork{
			FakeParameter: "yes",
			IWantFail:     false,
		},
	})
	assert.NoError(t, err)
	err = fake.CreateNetwork(nil, entity.Network{})
	assert.NoError(t, err)
}

func TestFakeNetworkValidateBeforeCreatingFail(t *testing.T) {
	fake, err := GetNetworkProvider(&entity.Network{
		Type: "fake",
		Fake: entity.FakeNetwork{
			FakeParameter: "",
		},
	})
	assert.NoError(t, err)
	err = fake.ValidateBeforeCreating(nil, entity.Network{})
	assert.Error(t, err)
}

func TestFakeNetworkCreatingFail(t *testing.T) {
	fake, err := GetNetworkProvider(&entity.Network{
		Type: "fake",
		Fake: entity.FakeNetwork{
			IWantFail: true,
		},
	})
	assert.NoError(t, err)
	err = fake.CreateNetwork(nil, entity.Network{})
	assert.Error(t, err)
}
