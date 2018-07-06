package storageprovider

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linkernetworks/vortex/src/entity"
)

func TestGetStorageProvider(t *testing.T) {
	testCases := []struct {
		cases               string
		storageType         entity.StorageType
		storageProviderType interface{}
	}{
		{"nfs", entity.NFSStorageType, reflect.TypeOf(NFSStorageProvider{})},
		{"fake", entity.FakeStorageType, reflect.TypeOf(FakeStorageProvider{})},
	}

	for _, tc := range testCases {
		t.Run(tc.cases, func(t *testing.T) {
			provider, err := GetStorageProvider(
				&entity.Storage{
					Type: tc.storageType,
					Fake: &entity.FakeStorage{},
					NFS:  &entity.NFSStorage{},
				})
			assert.NoError(t, err)
			a := reflect.TypeOf(provider)
			assert.Equal(t, a, tc.storageProviderType)
		})
	}
}

func TestGetStorageProviderFail(t *testing.T) {
	_, err := GetStorageProvider(
		&entity.Storage{
			Type: "Unknown",
		})
	assert.Error(t, err)
}
