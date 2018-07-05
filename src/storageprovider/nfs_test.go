package storageprovider

import (
	"bytes"
	"fmt"
	"math/rand"
	"os/exec"
	"runtime"
	"testing"
	"time"

	//"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	kc "github.com/linkernetworks/vortex/src/kubernetes"
	"github.com/linkernetworks/vortex/src/serviceprovider"
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
	if runtime.GOOS != "linux" {
		fmt.Println("We only testing the ovs function on Linux Host")
		t.Skip()
		return
	}
	suite.Run(t, new(StorageTestSuite))
}

func (suite *StorageTestSuite) TestCreateStorage() {
}
