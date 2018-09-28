package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/suite"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

type KubeCtlConfigMapTestSuite struct {
	suite.Suite
	kubectl    *KubeCtl
	fakeclient *fakeclientset.Clientset
}

func (suite *KubeCtlConfigMapTestSuite) SetupSuite() {
	suite.fakeclient = fakeclientset.NewSimpleClientset()
	suite.kubectl = New(suite.fakeclient)
}

func (suite *KubeCtlConfigMapTestSuite) TestGetConfigMap() {
	namespace := "default"
	configMap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "K8S-ConfigMap-1",
			Namespace: namespace,
		},
		Data: map[string]string{
			"firstData":  "first test content",
			"secondData": "second test content",
		},
	}
	_, err := suite.fakeclient.CoreV1().ConfigMaps(namespace).Create(&configMap)
	suite.NoError(err)

	result, err := suite.kubectl.GetConfigMap("K8S-ConfigMap-1", namespace)
	suite.NoError(err)
	suite.Equal(configMap.GetName(), result.GetName())
	suite.Equal(configMap.Data["firstData"], "first test content")
}

func (suite *KubeCtlConfigMapTestSuite) TestGetConfigMapFail() {
	namespace := "default"
	_, err := suite.kubectl.GetConfigMap("Unknown_Name", namespace)
	suite.Error(err)
}

func (suite *KubeCtlConfigMapTestSuite) TestGetConfigMaps() {
	namespace := "default"
	configMap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "K8S-ConfigMap-2",
			Namespace: namespace,
		},
		Data: map[string]string{
			"firstData":  "first test content",
			"secondData": "second test content",
		},
	}
	_, err := suite.fakeclient.CoreV1().ConfigMaps(namespace).Create(&configMap)
	suite.NoError(err)

	configMap = corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "K8S-ConfigMap-3",
			Namespace: namespace,
		},
		Data: map[string]string{
			"firstData":  "first test content",
			"secondData": "second test content",
		},
	}
	_, err = suite.fakeclient.CoreV1().ConfigMaps(namespace).Create(&configMap)
	suite.NoError(err)

	configMaps, err := suite.kubectl.GetConfigMaps(namespace)
	suite.NoError(err)
	suite.NotEqual(0, len(configMaps))
}

func (suite *KubeCtlConfigMapTestSuite) TestCreateDeleteConfigMap() {
	namespace := "default"
	configMap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "K8S-ConfigMap-4",
			Namespace: namespace,
		},
		Data: map[string]string{
			"firstData":  "first test content",
			"secondData": "second test content",
		},
	}
	_, err := suite.kubectl.CreateConfigMap(&configMap, namespace)
	suite.NoError(err)
	err = suite.kubectl.DeleteConfigMap("K8S-ConfigMap-4", namespace)
	suite.NoError(err)
}

func (suite *KubeCtlConfigMapTestSuite) TearDownSuite() {}

func TestKubeConfigMapTestSuite(t *testing.T) {
	suite.Run(t, new(KubeCtlConfigMapTestSuite))
}
