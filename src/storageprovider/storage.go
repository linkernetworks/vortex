package storageprovider

import (
	"fmt"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

// StorageProvider is storage provider interface
type StorageProvider interface {
	ValidateBeforeCreating(sp *serviceprovider.Container, net *entity.Storage) error
	CreateStorage(sp *serviceprovider.Container, net *entity.Storage) error
	DeleteStorage(sp *serviceprovider.Container, net *entity.Storage) error
}

// GetStorageProvider will get storage provider
func GetStorageProvider(storage *entity.Storage) (StorageProvider, error) {
	switch storage.Type {
	case "nfs":
		return NFSStorageProvider{*storage}, nil
	case "fake":
		return FakeStorageProvider{*storage.Fake}, nil
	default:
		return nil, fmt.Errorf("Unsupported Storage Type %s", storage.Type)
	}
}
