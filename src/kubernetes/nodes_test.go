package kubernetes

import (
	"testing"

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

func (suite *KubeCtlNodeTestSuite) SetupSuite() {
	suite.fakeclient = fakeclientset.NewSimpleClientset()
	suite.kubectl = New(suite.fakeclient)
}

func (suite *KubeCtlNodeTestSuite) TestGetNode() {
	node := corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Node-1",
		},
	}
	_, err := suite.fakeclient.CoreV1().Nodes().Create(&node)
	suite.NoError(err)

	result, err := suite.kubectl.GetNode("K8S-Node-1")
	suite.NoError(err)
	suite.Equal(node.GetName(), result.GetName())
}

func (suite *KubeCtlNodeTestSuite) TestGetNodeFail() {
	_, err := suite.kubectl.GetNode("Unknown_Name")
	suite.Error(err)
}

func (suite *KubeCtlNodeTestSuite) TestGetNodes() {
	node := corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Node-2",
		},
	}
	_, err := suite.fakeclient.CoreV1().Nodes().Create(&node)
	suite.NoError(err)

	node = corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Node-3",
		},
	}
	_, err = suite.fakeclient.CoreV1().Nodes().Create(&node)
	suite.NoError(err)

	nodes, err := suite.kubectl.GetNodes()
	suite.NoError(err)
	suite.NotEqual(0, len(nodes))
}

func (suite *KubeCtlNodeTestSuite) TestGetNodeExternalIP() {
	nodeAddr := corev1.NodeAddress{
		Type:    "ExternalIP",
		Address: "192.168.0.100",
	}
	node := corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Node-4",
		},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{nodeAddr},
		},
	}
	_, err := suite.fakeclient.CoreV1().Nodes().Create(&node)
	suite.NoError(err)

	nodeIP, err := suite.kubectl.GetNodeExternalIP("K8S-Node-4")
	suite.NoError(err)
	suite.Equal(nodeAddr.Address, nodeIP)
}

func (suite *KubeCtlNodeTestSuite) TestGetInvalidNodeExternalIP() {
	nodeAddr := corev1.NodeAddress{
		Type:    "ExternalIP",
		Address: "192.168.0.100",
	}
	node := corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Node-5",
		},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{nodeAddr},
		},
	}
	_, err := suite.fakeclient.CoreV1().Nodes().Create(&node)
	suite.NoError(err)

	nodeIP, err := suite.kubectl.GetNodeExternalIP("K8S-Node-99")
	suite.Error(err)
	suite.Equal("", nodeIP)
}

func (suite *KubeCtlNodeTestSuite) TestGetNodeInternalIP() {
	nodeAddr := corev1.NodeAddress{
		Type:    "InternalIP",
		Address: "10.0.2.200",
	}
	node := corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Node-6",
		},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{nodeAddr},
		},
	}
	_, err := suite.fakeclient.CoreV1().Nodes().Create(&node)
	suite.NoError(err)

	nodeIP, err := suite.kubectl.GetNodeInternalIP("K8S-Node-6")
	suite.NoError(err)
	suite.Equal(nodeAddr.Address, nodeIP)
}

func (suite *KubeCtlNodeTestSuite) TestGetInvalidNodeInternalIP() {
	nodeAddr := corev1.NodeAddress{
		Type:    "InternalIP",
		Address: "10.0.2.200",
	}
	node := corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "K8S-Node-7",
		},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{nodeAddr},
		},
	}
	_, err := suite.fakeclient.CoreV1().Nodes().Create(&node)
	suite.NoError(err)

	nodeIP, err := suite.kubectl.GetNodeInternalIP("K8S-Node-99")
	suite.Error(err)
	suite.Equal("", nodeIP)
}

func (suite *KubeCtlNodeTestSuite) TearDownSuite() {}

func TestKubeNodeTestSuite(t *testing.T) {
	suite.Run(t, new(KubeCtlNodeTestSuite))
}
