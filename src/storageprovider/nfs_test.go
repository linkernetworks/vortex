package storageprovider

import (
	"bytes"
	"fmt"
	"math/rand"
	"os/exec"
	"testing"
	"time"

	//"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	kc "github.com/linkernetworks/vortex/src/kubernetes"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	//"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2/bson"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func execute(suite *suite.Suite, cmd *exec.Cmd) {
	w := bytes.NewBuffer(nil)
	cmd.Stderr = w
	err := cmd.Run()
	suite.NoError(err)
	fmt.Printf("Stderr: %s\n", string(w.Bytes()))
}

type StorageTestSuite struct {
	suite.Suite
	sp *serviceprovider.Container
}

func (suite *StorageTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	suite.sp = serviceprovider.NewForTesting(cf)

	//init fakeclient
	fakeclient := fakeclientset.NewSimpleClientset()
	namespace := "default"
	suite.sp.KubeCtl = kc.New(fakeclient, namespace)

	//Create a fake clinet
}

func (suite *StorageTestSuite) TearDownSuite() {
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}

func (suite *StorageTestSuite) TestGetDeployment() {
	storage := &entity.Storage{
		Type: entity.NFSStorageType,
		NFS: entity.NFSStorage{
			IP:   "1.2.3.4",
			PATH: "/exports",
		},
	}

	deployment := getDeployment(bson.NewObjectId().Hex(), storage)
	suite.NotNil(deployment)
}

func (suite *StorageTestSuite) TestValidateBeforeCreating() {
	storage := &entity.Storage{
		Type: entity.NFSStorageType,
		NFS: entity.NFSStorage{
			IP:   "1.2.3.4",
			PATH: "/exports",
		},
	}

	//Parameters
	sp, err := GetStorageProvider(storage)
	suite.NoError(err)
	sp = sp.(NFSStorageProvider)

	err = sp.ValidateBeforeCreating(suite.sp, *storage)
	suite.NoError(err)
}

func (suite *StorageTestSuite) TestCreateStorage() {
	storage := entity.Storage{
		ID:   bson.NewObjectId(),
		Type: entity.NFSStorageType,
		NFS: entity.NFSStorage{
			IP:   "1.2.3.4",
			PATH: "/exports",
		},
	}

	sp, err := GetStorageProvider(&storage)
	suite.NoError(err)
	sp = sp.(NFSStorageProvider)

	err = sp.CreateStorage(suite.sp, storage)
	suite.NoError(err)

	deploy, err := suite.sp.KubeCtl.GetDeployment(NFS_PROVISIONER_PREFIX + storage.ID.Hex())
	suite.NotNil(deploy)
	suite.NoError(err)
}

func (suite *StorageTestSuite) TestDeleteStorage() {
	storage := entity.Storage{
		ID:   bson.NewObjectId(),
		Type: entity.NFSStorageType,
		NFS: entity.NFSStorage{
			IP:   "1.2.3.4",
			PATH: "/exports",
		},
	}

	sp, err := GetStorageProvider(&storage)
	suite.NoError(err)
	sp = sp.(NFSStorageProvider)

	err = sp.CreateStorage(suite.sp, storage)
	suite.NoError(err)

	deploy, err := suite.sp.KubeCtl.GetDeployment(NFS_PROVISIONER_PREFIX + storage.ID.Hex())
	suite.NotNil(deploy)
	suite.NoError(err)

	err = sp.DeleteStorage(suite.sp, storage)
	suite.NoError(err)

	deploy, err = suite.sp.KubeCtl.GetDeployment(NFS_PROVISIONER_PREFIX + storage.ID.Hex())
	suite.Nil(deploy)
	suite.Error(err)
}

func (suite *StorageTestSuite) TestValidateBeforeCreatingFail() {
	testCases := []struct {
		caseName string
		storage  *entity.Storage
	}{
		{"invalidIP", &entity.Storage{
			Type: entity.NFSStorageType,
			NFS: entity.NFSStorage{
				IP: "a.b.c.d",
			},
		}},
		{"invalidExports-1", &entity.Storage{
			Type: entity.NFSStorageType,
			NFS: entity.NFSStorage{
				IP:   "1.2.3.4",
				PATH: "tmp",
			},
		}},
		{"invalidExports-2", &entity.Storage{
			Type: entity.NFSStorageType,
			NFS: entity.NFSStorage{
				IP:   "1.2.3.4",
				PATH: "",
			},
		}},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.caseName, func(t *testing.T) {
			//Parameters
			np, err := GetStorageProvider(tc.storage)
			suite.NoError(err)
			np = np.(NFSStorageProvider)

			err = np.ValidateBeforeCreating(suite.sp, *tc.storage)
			suite.Error(err)
		})
	}
}
