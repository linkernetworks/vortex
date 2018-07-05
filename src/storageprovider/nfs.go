package storageprovider

import (
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	//	"gopkg.in/mgo.v2/bson"
)

type NFSStorageProvider struct {
	entity.NFSStorage
}

func (nfs NFSStorageProvider) ValidateBeforeCreating(sp *serviceprovider.Container, network entity.Storage) error {
	return nil
}

func (nfs NFSStorageProvider) CreateStorage(sp *serviceprovider.Container, network entity.Storage) error {
	return nil
}

func (nfs NFSStorageProvider) DeleteStorage(sp *serviceprovider.Container, network entity.Storage) error {
	return nil
}
