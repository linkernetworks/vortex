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
		{"system", entity.OVSKernelspaceNetworkType, reflect.TypeOf(OVSNetworkProvider{})},
		{"netdev", entity.OVSUserspaceNetworkType, reflect.TypeOf(OVSUserspaceNetworkProvider{})},
		{"fake", entity.FakeNetworkType, reflect.TypeOf(FakeNetworkProvider{})},
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
