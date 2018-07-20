package networkprovider

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linkernetworks/vortex/src/entity"
)

func TestGetNetworkProvider(t *testing.T) {
	testCases := []struct {
		cases           string
		netType         entity.NetworkType
		netProviderType interface{}
	}{
		{"system", entity.OVSKernelspaceNetworkType, reflect.TypeOf(kernelspaceNetworkProvider{})},
		{"netdev", entity.OVSUserspaceNetworkType, reflect.TypeOf(userspaceNetworkProvider{})},
		{"fake", entity.FakeNetworkType, reflect.TypeOf(fakeNetworkProvider{})},
	}

	for _, tc := range testCases {
		t.Run(tc.cases, func(t *testing.T) {
			provider, _ := GetNetworkProvider(
				&entity.Network{
					Type: tc.netType,
				})
			a := reflect.TypeOf(provider)
			assert.Equal(t, a, tc.netProviderType)
		})
	}
}

func TestGetNetworkProviderFail(t *testing.T) {
	_, err := GetNetworkProvider(
		&entity.Network{
			Type: "Unknown",
		})
	assert.Error(t, err)
}

func TestGenerateBridgeName(t *testing.T) {
	ans := GenerateBridgeName("netdev", "my network 1")
	assert.Equal(t, "netdev-de0165", ans)
}
