package storageprovider

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linkernetworks/vortex/src/entity"
)

func TestStorageValidateBeforeCreating(t *testing.T) {
	fake, err := GetStorageProvider(&entity.Storage{
		Type: "fake",
		Fake: &entity.FakeStorage{
			FakeParameter: "yes",
		},
	})
	assert.NoError(t, err)
	err = fake.ValidateBeforeCreating(nil, &entity.Storage{})
	assert.NoError(t, err)
}

func TestFakeStorageCreating(t *testing.T) {
	fake, err := GetStorageProvider(&entity.Storage{
		Type: "fake",
		Fake: &entity.FakeStorage{
			FakeParameter: "yes",
			IWantFail:     false,
		},
	})
	assert.NoError(t, err)
	err = fake.CreateStorage(nil, &entity.Storage{})
	assert.NoError(t, err)
}

func TestFakeStorageValidateBeforeCreatingFail(t *testing.T) {
	fake, err := GetStorageProvider(&entity.Storage{
		Type: "fake",
		Fake: &entity.FakeStorage{
			FakeParameter: "",
		},
	})
	assert.NoError(t, err)
	err = fake.ValidateBeforeCreating(nil, &entity.Storage{})
	assert.Error(t, err)
}

func TestFakeStorageCreatingFail(t *testing.T) {
	fake, err := GetStorageProvider(&entity.Storage{
		Type: "fake",
		Fake: &entity.FakeStorage{
			IWantFail: true,
		},
	})
	assert.NoError(t, err)
	err = fake.CreateStorage(nil, &entity.Storage{})
	assert.Error(t, err)
}

func TestFakeStorageDelete(t *testing.T) {
	fake, err := GetStorageProvider(&entity.Storage{
		Type: "fake",
		Fake: &entity.FakeStorage{},
	})
	assert.NoError(t, err)
	err = fake.DeleteStorage(nil, &entity.Storage{})
	assert.NoError(t, err)
}

func TestFakeStorageDeleteFail(t *testing.T) {
	fake, err := GetStorageProvider(&entity.Storage{
		Type: "fake",
		Fake: &entity.FakeStorage{
			IWantFail: true,
		},
	})
	assert.NoError(t, err)
	err = fake.DeleteStorage(nil, &entity.Storage{})
	assert.Error(t, err)
}
