package volume

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

type VolumeTestSuite struct {
	suite.Suite
	sp *serviceprovider.Container
}

func (suite *VolumeTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	suite.sp = serviceprovider.NewForTesting(cf)
}

func (suite *VolumeTestSuite) TearDownSuite() {
}

func TestVolumeSuite(t *testing.T) {
	suite.Run(t, new(VolumeTestSuite))
}

func (suite *VolumeTestSuite) TestGetPVCInstance() {
	volume := &entity.Volume{
		ID:          bson.NewObjectId(),
		Name:        namesgenerator.GetRandomName(0),
		StorageName: namesgenerator.GetRandomName(0),
	}

	pvc := getPVCInstance(volume, namesgenerator.GetRandomName(0), namesgenerator.GetRandomName(0))
	suite.NotNil(pvc)
}

func (suite *VolumeTestSuite) TestGetStorageClassName() {
	session := suite.sp.Mongo.NewSession()
	defer session.Close()

	storage := entity.Storage{
		ID:               bson.NewObjectId(),
		Name:             namesgenerator.GetRandomName(0),
		StorageClassName: namesgenerator.GetRandomName(0),
	}

	session.Insert(entity.StorageCollectionName, &storage)
	defer session.Remove(entity.StorageCollectionName, "name", storage.Name)
	name, err := getStorageClassName(session, storage.Name)
	suite.NoError(err)
	suite.Equal(name, storage.StorageClassName)
}

func (suite *VolumeTestSuite) TestCreateVolume() {
	session := suite.sp.Mongo.NewSession()
	defer session.Close()
	storage := entity.Storage{
		ID:               bson.NewObjectId(),
		Name:             namesgenerator.GetRandomName(0),
		StorageClassName: namesgenerator.GetRandomName(0),
	}

	session.Insert(entity.StorageCollectionName, &storage)
	defer session.Remove(entity.StorageCollectionName, "name", storage.Name)

	volume := &entity.Volume{
		ID:          bson.NewObjectId(),
		Name:        namesgenerator.GetRandomName(0),
		StorageName: storage.Name,
	}

	err := CreateVolume(suite.sp, volume)
	suite.NoError(err)

	name := volume.GetPVCName()
	v, err := suite.sp.KubeCtl.GetPVC(name)
	suite.NoError(err)
	suite.NotNil(v)

	err = DeleteVolume(suite.sp, volume)
	suite.NoError(err)

	v, err = suite.sp.KubeCtl.GetPVC(name)
	suite.Error(err)
	suite.Nil(v)
}

func (suite *VolumeTestSuite) TestCreateVolumeFail() {
	volume := &entity.Volume{
		ID:          bson.NewObjectId(),
		Name:        namesgenerator.GetRandomName(0),
		StorageName: namesgenerator.GetRandomName(0),
	}

	err := CreateVolume(suite.sp, volume)
	suite.Error(err)
}
