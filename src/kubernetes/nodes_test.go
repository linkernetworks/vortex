package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

type KubeCtlTestSuite struct {
	suite.Suite
	kubectl    *KubeCtl
	fakeclient fakeclientset.Clientset
}

func (suite *KubeCtlTestSuite) SetupTest() {
	clientset := fakeclientset.NewSimpleClientset()
	namespace := "default"
	suite.kubectl = New(clientset, namespace)
}

func (suite *KubeCtlTestSuite) TestGetNode(t *testing.T) {
	node := corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Node",
		},
	}
	_, err := suite.fakeclient.CoreV1().Nodes().Create(&node)
	assert.NoError(t, err)

	result, err := suite.kubectl.GetNode("K8S-Node")
	assert.NoError(t, err)
	assert.Equal(t, node.GetName(), result.GetName())
}

func (suite *KubeCtlTestSuite) TestGetNodeFail(t *testing.T) {
	_, err := suite.kubectl.GetNode("UnKnown_Name")
	assert.Error(t, err)
}

func (suite *KubeCtlTestSuite) TestGetNodes(t *testing.T) {
	node := corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Node1",
		},
	}
	_, err := suite.fakeclient.CoreV1().Nodes().Create(&node)
	assert.NoError(t, err)

	node = corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Node2",
		},
	}
	_, err = suite.fakeclient.CoreV1().Nodes().Create(&node)
	assert.NoError(t, err)

	nodes, err := suite.kubectl.GetNodes()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(nodes))
}

func (suite *KubeCtlTestSuite) TestGetNodeExternalIP(t *testing.T) {
	nodeAddr := corev1.NodeAddress{
		Type:    "ExternalIP",
		Address: "192.168.0.100",
	}
	node := corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Node",
		},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{nodeAddr},
		},
	}
	_, err := suite.fakeclient.CoreV1().Nodes().Create(&node)
	assert.NoError(t, err)

	nodeIP, err := suite.kubectl.GetNodeExternalIP("K8S-Node")
	assert.NoError(t, err)
	assert.Equal(t, nodeAddr.Address, nodeIP)
}
