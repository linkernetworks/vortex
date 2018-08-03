package storageprovider

import (
	"math/rand"
	"testing"
	"time"

	//"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"

	"gopkg.in/mgo.v2/bson"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type StorageTestSuite struct {
	suite.Suite
	sp *serviceprovider.Container
}

func (suite *StorageTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	suite.sp = serviceprovider.NewForTesting(cf)
}

func (suite *StorageTestSuite) TearDownSuite() {
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}

func (suite *StorageTestSuite) TestGetDeployment() {
	storage := &entity.Storage{
		Type: entity.NFSStorageType,
		IP:   "1.2.3.4",
		PATH: "/exports",
	}

	deployment := getDeployment(bson.NewObjectId().Hex(), storage)
	suite.NotNil(deployment)
}

func (suite *StorageTestSuite) TestGetStorageClass() {
	storage := &entity.Storage{
		Type: entity.NFSStorageType,
		IP:   "1.2.3.4",
		PATH: "/exports",
	}

	storageClass := getStorageClass(bson.NewObjectId().Hex(), bson.NewObjectId().Hex(), storage)
	suite.NotNil(storageClass)
}

func (suite *StorageTestSuite) TestValidateBeforeCreating() {
	storage := &entity.Storage{
		Type: entity.NFSStorageType,
		IP:   "1.2.3.4",
		PATH: "/exports",
	}

	//Parameters
	sp, err := GetStorageProvider(storage)
	suite.NoError(err)
	sp = sp.(NFSStorageProvider)

	err = sp.ValidateBeforeCreating(suite.sp, storage)
	suite.NoError(err)
}

func (suite *StorageTestSuite) TestCreateStorage() {
	namespace := "vortex"
	storage := entity.Storage{
		ID:   bson.NewObjectId(),
		Type: entity.NFSStorageType,
		IP:   "1.2.3.4",
		PATH: "/exports",
	}

	sp, err := GetStorageProvider(&storage)
	suite.NoError(err)
	sp = sp.(NFSStorageProvider)

	err = sp.CreateStorage(suite.sp, &storage)
	suite.NoError(err)

	deploy, err := suite.sp.KubeCtl.GetDeployment(NFSProvisionerPrefix+storage.ID.Hex(), namespace)
	suite.NotNil(deploy)
	suite.NoError(err)
}

func (suite *StorageTestSuite) TestDeleteStorage() {
	namespace := "vortex"
	storage := &entity.Storage{
		ID:   bson.NewObjectId(),
		Type: entity.NFSStorageType,
		IP:   "1.2.3.4",
		PATH: "/exports",
	}

	sp, err := GetStorageProvider(storage)
	suite.NoError(err)
	sp = sp.(NFSStorageProvider)

	err = sp.CreateStorage(suite.sp, storage)
	suite.NoError(err)

	deploy, err := suite.sp.KubeCtl.GetDeployment(NFSProvisionerPrefix+storage.ID.Hex(), namespace)
	suite.NotNil(deploy)
	suite.NoError(err)

	err = sp.DeleteStorage(suite.sp, storage)
	suite.NoError(err)

	deploy, err = suite.sp.KubeCtl.GetDeployment(NFSProvisionerPrefix+storage.ID.Hex(), namespace)
	suite.Nil(deploy)
	suite.Error(err)
}

func (suite *StorageTestSuite) TestDeleteStorageFail() {
	storage := &entity.Storage{
		ID:   bson.NewObjectId(),
		Type: entity.NFSStorageType,
		IP:   "1.2.3.4",
		PATH: "/exports",
		Name: namesgenerator.GetRandomName(0),
	}

	volume := entity.Volume{
		ID:          bson.NewObjectId(),
		StorageName: storage.Name,
	}

	session := suite.sp.Mongo.NewSession()
	defer session.Close()

	session.Insert(entity.VolumeCollectionName, volume)
	defer session.Remove(entity.VolumeCollectionName, "_id", volume.ID)

	sp, err := GetStorageProvider(storage)
	suite.NoError(err)
	sp = sp.(NFSStorageProvider)

	err = sp.DeleteStorage(suite.sp, storage)
	suite.Error(err)
}

func (suite *StorageTestSuite) TestValidateBeforeCreatingFail() {
	testCases := []struct {
		caseName string
		storage  *entity.Storage
	}{
		{"invalidIP", &entity.Storage{
			Type: entity.NFSStorageType,
			IP:   "a.b.c.d",
		}},
		{"invalidExports-1", &entity.Storage{
			Type: entity.NFSStorageType,
			IP:   "1.2.3.4",
			PATH: "tmp",
		}},
		{"invalidExports-2", &entity.Storage{
			Type: entity.NFSStorageType,
			IP:   "1.2.3.4",
			PATH: "",
		}},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.caseName, func(t *testing.T) {
			//Parameters
			np, err := GetStorageProvider(tc.storage)
			suite.NoError(err)
			np = np.(NFSStorageProvider)

			err = np.ValidateBeforeCreating(suite.sp, tc.storage)
			suite.Error(err)
		})
	}
}
