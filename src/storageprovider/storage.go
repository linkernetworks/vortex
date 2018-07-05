package storageprovider

import (
	"fmt"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

type StorageProvider interface {
	ValidateBeforeCreating(sp *serviceprovider.Container, net *entity.Storage) error
	CreateStorage(sp *serviceprovider.Container, net *entity.Storage) error
	DeleteStorage(sp *serviceprovider.Container, net *entity.Storage) error
}

func GetStorageProvider(storage *entity.Storage) (StorageProvider, error) {
	switch storage.Type {
	case "nfs":
		return NFSStorageProvider{storage.NFS}, nil
	case "fake":
		return FakeStorageProvider{storage.Fake}, nil
	default:
		return nil, fmt.Errorf("Unsupported Storage Type %s", storage.Type)
	}
}
