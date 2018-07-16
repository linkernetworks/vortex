package kubernetes

import (
	"math/rand"
	"testing"
	"time"

	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/storage/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

type KubeCtlStorageClassTestSuite struct {
	suite.Suite
	kubectl    *KubeCtl
	fakeclient *fakeclientset.Clientset
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (suite *KubeCtlStorageClassTestSuite) SetupSuite() {
	suite.fakeclient = fakeclientset.NewSimpleClientset()
	suite.kubectl = New(suite.fakeclient)
}

func (suite *KubeCtlStorageClassTestSuite) TearDownSuite() {}
func (suite *KubeCtlStorageClassTestSuite) TestCreateStorageClass() {
	name := namesgenerator.GetRandomName(0)
	storageClass := v1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Provisioner: "",
	}
	ret, err := suite.kubectl.CreateStorageClass(&storageClass)
	suite.NoError(err)
	suite.NotNil(ret)

	deploy, err := suite.kubectl.GetStorageClass(name)
	suite.NoError(err)
	suite.NotNil(deploy)
}

func (suite *KubeCtlStorageClassTestSuite) TestDeleteStorageClass() {
	name := namesgenerator.GetRandomName(0)
	storageClass := v1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Provisioner: "",
	}
	ret, err := suite.kubectl.CreateStorageClass(&storageClass)
	suite.NoError(err)
	suite.NotNil(ret)

	deploy, err := suite.kubectl.GetStorageClass(name)
	suite.NoError(err)
	suite.NotNil(deploy)

	err = suite.kubectl.DeleteStorageClass(name)
	suite.NoError(err)
	deploy, err = suite.kubectl.GetStorageClass(name)
	suite.Error(err)
	suite.Nil(deploy)
}

func TestStorageClassTestSuite(t *testing.T) {
	suite.Run(t, new(KubeCtlStorageClassTestSuite))
}
