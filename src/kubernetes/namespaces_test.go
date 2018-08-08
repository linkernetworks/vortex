package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/suite"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

type KubeCtlNamespaceTestSuite struct {
	suite.Suite
	kubectl    *KubeCtl
	fakeclient *fakeclientset.Clientset
}

func (suite *KubeCtlNamespaceTestSuite) SetupSuite() {
	suite.fakeclient = fakeclientset.NewSimpleClientset()
	suite.kubectl = New(suite.fakeclient)
}

func (suite *KubeCtlNamespaceTestSuite) TestGetNamespace() {
	namespace := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Namespace-1",
		},
	}
	_, err := suite.fakeclient.CoreV1().Namespaces().Create(&namespace)
	suite.NoError(err)

	result, err := suite.kubectl.GetNamespace("K8S-Namespace-1")
	suite.NoError(err)
	suite.Equal(namespace.GetName(), result.GetName())
}

func (suite *KubeCtlNamespaceTestSuite) TestGetNamespaceFail() {
	_, err := suite.kubectl.GetNamespace("Unknown_Name")
	suite.Error(err)
}

func (suite *KubeCtlNamespaceTestSuite) TestGetNamespaces() {
	namespace := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Namespace-2",
		},
	}
	_, err := suite.fakeclient.CoreV1().Namespaces().Create(&namespace)
	suite.NoError(err)

	namespace = corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Namespace-3",
		},
	}
	_, err = suite.fakeclient.CoreV1().Namespaces().Create(&namespace)
	suite.NoError(err)

	namespaces, err := suite.kubectl.GetNamespaces()
	suite.NoError(err)
	suite.NotEqual(0, len(namespaces))
}

func (suite *KubeCtlNamespaceTestSuite) TestCreateDeleteNamespace() {
	namespace := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Namespace-4",
		},
	}
	_, err := suite.kubectl.CreateNamespace(&namespace)
	suite.NoError(err)
	err = suite.kubectl.DeleteNamespace("K8S-Namespace-4")
	suite.NoError(err)
}

func (suite *KubeCtlNamespaceTestSuite) TearDownSuite() {}

func TestKubeNamespaceTestSuite(t *testing.T) {
	suite.Run(t, new(KubeCtlNamespaceTestSuite))
}
