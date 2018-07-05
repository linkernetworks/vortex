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

	//mgo "gopkg.in/mgo.v2"

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
	sp       *serviceprovider.Container
	fakeName string //Use for non-connectivity node
	np       NFSStorageProvider
	network  entity.Storage
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

func (suite *StorageTestSuite) TestCreateStorage() {

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
