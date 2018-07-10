package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/suite"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

type KubeCtlServiceTestSuite struct {
	suite.Suite
	kubectl    *KubeCtl
	fakeclient *fakeclientset.Clientset
}

func (suite *KubeCtlServiceTestSuite) SetupSuite() {
	suite.fakeclient = fakeclientset.NewSimpleClientset()
	namespace := "default"
	suite.kubectl = New(suite.fakeclient, namespace)
}

func (suite *KubeCtlServiceTestSuite) TestGetService() {
	namespace := "default"
	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Service-1",
		},
	}
	_, err := suite.fakeclient.CoreV1().Services(namespace).Create(&service)
	suite.NoError(err)

	result, err := suite.kubectl.GetService("K8S-Service-1")
	suite.NoError(err)
	suite.Equal(service.GetName(), result.GetName())
}

func (suite *KubeCtlServiceTestSuite) TestGetServiceFail() {
	_, err := suite.kubectl.GetService("Unknown_Name")
	suite.Error(err)
}

func (suite *KubeCtlServiceTestSuite) TestGetServices() {
	namespace := "default"
	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Service-2",
		},
	}
	_, err := suite.fakeclient.CoreV1().Services(namespace).Create(&service)
	suite.NoError(err)

	service = corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Service-3",
		},
	}
	_, err = suite.fakeclient.CoreV1().Services(namespace).Create(&service)
	suite.NoError(err)

	services, err := suite.kubectl.GetServices()
	suite.NoError(err)
	suite.NotEqual(0, len(services))
}

func (suite *KubeCtlServiceTestSuite) TestCreateDeleteService() {
	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Service-4",
		},
	}
	_, err := suite.kubectl.CreateService(&service)
	suite.NoError(err)
	err = suite.kubectl.DeleteService("K8S-Service-4")
	suite.NoError(err)
}

func (suite *KubeCtlServiceTestSuite) TearDownSuite() {}

func TestKubeServiceTestSuite(t *testing.T) {
	suite.Run(t, new(KubeCtlServiceTestSuite))
}
