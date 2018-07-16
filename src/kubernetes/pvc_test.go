package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/suite"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

type KubeCtlPVCTestSuite struct {
	suite.Suite
	kubectl    *KubeCtl
	fakeclient *fakeclientset.Clientset
}

func (suite *KubeCtlPVCTestSuite) SetupSuite() {
	namespace := "default"
	suite.fakeclient = fakeclientset.NewSimpleClientset()
	suite.kubectl = New(suite.fakeclient, namespace)
}

func (suite *KubeCtlPVCTestSuite) TestGetPVC() {
	namespace := "default"
	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-PVC-1",
		},
	}
	_, err := suite.fakeclient.CoreV1().PersistentVolumeClaims(namespace).Create(&pvc)
	suite.NoError(err)

	result, err := suite.kubectl.GetPVC("K8S-PVC-1", namespace)
	suite.NoError(err)
	suite.Equal(pvc.GetName(), result.GetName())
}

func (suite *KubeCtlPVCTestSuite) TestGetPVCFail() {
	namespace := "default"
	_, err := suite.kubectl.GetPVC("Unknown_Name", namespace)
	suite.Error(err)
}

func (suite *KubeCtlPVCTestSuite) TestGetPVCs() {
	namespace := "default"
	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-PVC-2",
		},
	}
	_, err := suite.fakeclient.CoreV1().PersistentVolumeClaims(namespace).Create(&pvc)
	suite.NoError(err)

	pvc = corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-PVC-3",
		},
	}
	_, err = suite.fakeclient.CoreV1().PersistentVolumeClaims(namespace).Create(&pvc)
	suite.NoError(err)

	pvcs, err := suite.kubectl.GetPVCs(namespace)
	suite.NoError(err)
	suite.NotEqual(0, len(pvcs))
}

func (suite *KubeCtlPVCTestSuite) TestCreateDeletePVC() {
	namespace := "default"
	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-PVC-4",
		},
	}
	_, err := suite.kubectl.CreatePVC(&pvc, namespace)
	suite.NoError(err)
	err = suite.kubectl.DeletePVC("K8S-PVC-4", namespace)
	suite.NoError(err)
}

func (suite *KubeCtlPVCTestSuite) TearDownSuite() {}

func TestKubePVCTestSuite(t *testing.T) {
	suite.Run(t, new(KubeCtlPVCTestSuite))
}
