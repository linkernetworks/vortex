package storageprovider

import (
	"fmt"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"net"
	//	"gopkg.in/mgo.v2/bson"
)

type NFSStorageProvider struct {
	entity.NFSStorage
}

func (nfs NFSStorageProvider) ValidateBeforeCreating(sp *serviceprovider.Container, storage entity.Storage) error {
	ip := net.ParseIP(storage.NFS.IP)
	if len(ip) == 0 {
		return fmt.Errorf("Invalid IP address %s\n", storage.NFS.IP)
	}

	path := storage.NFS.PATH
	if path == "" || path[0] != '/' {
		return fmt.Errorf("Invalid NFS export path %s\n", path)
	}

	return nil
}

func (nfs NFSStorageProvider) CreateStorage(sp *serviceprovider.Container, storage entity.Storage) error {
	return nil
}

func (nfs NFSStorageProvider) DeleteStorage(sp *serviceprovider.Container, storage entity.Storage) error {
	return nil
}
