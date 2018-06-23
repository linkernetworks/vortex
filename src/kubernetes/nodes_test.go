package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

type KubeCtlNodeTestSuite struct {
	suite.Suite
	kubectl    *KubeCtl
	fakeclient *fakeclientset.Clientset
}

func (suite *KubeCtlNodeTestSuite) SetupTest() {
	suite.fakeclient = fakeclientset.NewSimpleClientset()
	namespace := "default"
	suite.kubectl = New(suite.fakeclient, namespace)
}

func (suite *KubeCtlNodeTestSuite) TestGetNode() {
	node := corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Node",
		},
	}
	_, err := suite.fakeclient.CoreV1().Nodes().Create(&node)
	assert.NoError(suite.T(), err)

	result, err := suite.kubectl.GetNode("K8S-Node")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), node.GetName(), result.GetName())
}

func (suite *KubeCtlNodeTestSuite) TestGetNodeFail() {
	_, err := suite.kubectl.GetNode("Unknown_Name")
	assert.Error(suite.T(), err)
}

func (suite *KubeCtlNodeTestSuite) TestGetNodes() {
	node := corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Node1",
		},
	}
	_, err := suite.fakeclient.CoreV1().Nodes().Create(&node)
	assert.NoError(suite.T(), err)

	node = corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Node2",
		},
	}
	_, err = suite.fakeclient.CoreV1().Nodes().Create(&node)
	assert.NoError(suite.T(), err)

	nodes, err := suite.kubectl.GetNodes()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(nodes))
}

func (suite *KubeCtlNodeTestSuite) TestGetNodeExternalIP() {
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
	assert.NoError(suite.T(), err)

	nodeIP, err := suite.kubectl.GetNodeExternalIP("K8S-Node")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), nodeAddr.Address, nodeIP)
}

func (suite *KubeCtlNodeTestSuite) TestGetInvalidNodeExternalIP() {
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
	assert.NoError(suite.T(), err)

	nodeIP, err := suite.kubectl.GetNodeExternalIP("K8S-Node-0")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "", nodeIP)
}

func (suite *KubeCtlNodeTestSuite) TearDownTest() {}

func TestKubeNodeTestSuite(t *testing.T) {
	suite.Run(t, new(KubeCtlNodeTestSuite))
}
