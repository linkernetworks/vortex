package kubernetes

import (
	"math/rand"
	"testing"
	"time"

	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"
	appsv1 "k8s.io/api/apps/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

type KubeCtlDeploymentTestSuite struct {
	suite.Suite
	kubectl    *KubeCtl
	fakeclient *fakeclientset.Clientset
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (suite *KubeCtlDeploymentTestSuite) SetupSuite() {
	suite.fakeclient = fakeclientset.NewSimpleClientset()
	namespace := "default"
	suite.kubectl = New(suite.fakeclient, namespace)
}

func (suite *KubeCtlDeploymentTestSuite) TearDownSuite() {}
func (suite *KubeCtlDeploymentTestSuite) TestCreateDeployment() {
	var replicas int32
	replicas = 3
	name := namesgenerator.GetRandomName(0)
	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
		},
		Status: appsv1.DeploymentStatus{},
	}
	ret, err := suite.kubectl.CreateDeployment(&deployment)
	suite.NoError(err)
	suite.NotNil(ret)

	deploy, err := suite.kubectl.GetDeployment(name)
	suite.NoError(err)
	suite.NotNil(deploy)
	suite.Equal(replicas, *deploy.Spec.Replicas)
}

func (suite *KubeCtlDeploymentTestSuite) TestDeleteDeployment() {
	var replicas int32
	replicas = 3
	name := namesgenerator.GetRandomName(0)
	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
		},
		Status: appsv1.DeploymentStatus{},
	}
	ret, err := suite.kubectl.CreateDeployment(&deployment)
	suite.NoError(err)
	suite.NotNil(ret)

	deploy, err := suite.kubectl.GetDeployment(name)
	suite.NoError(err)
	suite.NotNil(deploy)
	suite.Equal(replicas, *deploy.Spec.Replicas)

	err = suite.kubectl.DeleteDeployment(name)
	suite.NoError(err)
	deploy, err = suite.kubectl.GetDeployment(name)
	suite.Error(err)
	suite.Nil(deploy)
}

func TestDeploymentTestSuite(t *testing.T) {
	suite.Run(t, new(KubeCtlDeploymentTestSuite))
}
