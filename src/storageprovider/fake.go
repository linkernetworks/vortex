package storageprovider

import (
	"fmt"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

type FakeStorageProvider struct {
	entity.FakeStorage
}

func (fake FakeStorageProvider) ValidateBeforeCreating(sp *serviceprovider.Container, net entity.Storage) error {
	if fake.FakeParameter == "" {
		return fmt.Errorf("Fail to validate but don't worry, I'm fake storage provider")
	}
	return nil
}

func (fake FakeStorageProvider) CreateStorage(sp *serviceprovider.Container, net entity.Storage) error {
	if fake.IWantFail {
		return fmt.Errorf("Fail to create storage but don't worry, I'm fake storage provider")
	}
	return nil
}

func (fake FakeStorageProvider) DeleteStorage(sp *serviceprovider.Container, net entity.Storage) error {
	if fake.IWantFail {
		return fmt.Errorf("Fail to delete storage but don't worry, I'm fake storage provider")
	}
	return nil
}
