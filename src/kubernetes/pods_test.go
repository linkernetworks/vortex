package kubernetes

import (
	"math/rand"
	"testing"
	"time"

	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type KubeCtlPodTestSuite struct {
	suite.Suite
	kubectl    *KubeCtl
	fakeclient *fakeclientset.Clientset
}

func (suite *KubeCtlPodTestSuite) SetupSuite() {
	suite.fakeclient = fakeclientset.NewSimpleClientset()
	suite.kubectl = New(suite.fakeclient)
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
	namespace := "default"
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

func (suite *KubeCtlPodTestSuite) TestDoesPodCompleted() {
	namespace := "default"
	pods := []corev1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: namesgenerator.GetRandomName(0),
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodPending,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: namesgenerator.GetRandomName(0),
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodSucceeded,
			},
		},
	}

	for _, pod := range pods {
		_, err := suite.kubectl.CreatePod(&pod, namespace)
		suite.NoError(err)
	}

	run := suite.kubectl.IsPodCompleted(&pods[0])
	suite.False(run)

	run = suite.kubectl.IsPodCompleted(&pods[1])
	suite.True(run)
}

func (suite *KubeCtlPodTestSuite) TearDownSuite() {}

func TestKubePodTestSuite(t *testing.T) {
	suite.Run(t, new(KubeCtlPodTestSuite))
}
