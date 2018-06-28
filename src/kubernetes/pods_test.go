package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/suite"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

type KubeCtlPodTestSuite struct {
	suite.Suite
	kubectl    *KubeCtl
	fakeclient *fakeclientset.Clientset
}

func (suite *KubeCtlPodTestSuite) SetupSuite() {
	suite.fakeclient = fakeclientset.NewSimpleClientset()
	namespace := "default"
	suite.kubectl = New(suite.fakeclient, namespace)
}

func (suite *KubeCtlPodTestSuite) TestGetPod() {
	namespace := "default"
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Pod-1",
		},
	}
	_, err := suite.fakeclient.CoreV1().Pods(namespace).Create(&pod)
	suite.NoError(err)

	result, err := suite.kubectl.GetPod("K8S-Pod-1", namespace)
	suite.NoError(err)
	suite.Equal(pod.GetName(), result.GetName())
}

func (suite *KubeCtlPodTestSuite) TestGetPodFail() {
	namespace := "Unknown_Namespace"
	_, err := suite.kubectl.GetPod("Unknown_Name", namespace)
	suite.Error(err)
}

func (suite *KubeCtlPodTestSuite) TestGetPods() {
	namespace := "default"
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Pod-2",
		},
	}
	_, err := suite.fakeclient.CoreV1().Pods(namespace).Create(&pod)
	suite.NoError(err)

	pod = corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Pod-3",
		},
	}
	_, err = suite.fakeclient.CoreV1().Pods(namespace).Create(&pod)
	suite.NoError(err)

	pods, err := suite.kubectl.GetPods(namespace)
	suite.NoError(err)
	suite.NotEqual(0, len(pods))
}

func (suite *KubeCtlPodTestSuite) TestCreateDeletePod() {
	namespace := "default"
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Pod-4",
		},
	}
	_, err := suite.kubectl.CreatePod(&pod, namespace)
	suite.NoError(err)
	err = suite.kubectl.DeletePod("K8S-Pod-4", namespace)
	suite.NoError(err)
}

func (suite *KubeCtlPodTestSuite) TearDownSuite() {}

func TestKubePodTestSuite(t *testing.T) {
	suite.Run(t, new(KubeCtlPodTestSuite))
}
