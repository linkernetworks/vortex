package networkcontroller

import (
	"fmt"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/kubernetes"
	"github.com/stretchr/testify/suite"
	"os"
	"os/exec"
	"runtime"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

type NetworkControllerTestSuite struct {
	suite.Suite
	kubectl  *kubernetes.KubeCtl
	ifName   string
	nodeName string
}

func (suite *NetworkControllerTestSuite) SetupSuite() {
	// init fakeclient
	fakeclient := fakeclientset.NewSimpleClientset()
	namespace := "default"
	suite.kubectl = kubernetes.New(fakeclient, namespace)

	//Create a fake clinet
	//Init
	nodeAddr := corev1.NodeAddress{
		Type:    "ExternalIP",
		Address: "127.0.0.1",
	}

	suite.nodeName = namesgenerator.GetRandomName(0)[0:8]
	node := corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: suite.nodeName,
		},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{nodeAddr},
		},
	}
	_, err := suite.kubectl.Clientset.CoreV1().Nodes().Create(&node)
	suite.NoError(err)

	//There's a length limit of link name
	suite.ifName = namesgenerator.GetRandomName(0)[0:8]
	pName := namesgenerator.GetRandomName(0)[0:8]
	//Create a veth for testing
	err = exec.Command("ip", "link", "add", suite.ifName, "type", "veth", "peer", "name", pName).Run()
	suite.NoError(err)
}

func TestNetworkControllerSuite(t *testing.T) {
	if runtime.GOOS != "linux" {
		fmt.Println("We only testing the ovs function on Linux Host")
		t.Skip()
		return
	}
	if _, defined := os.LookupEnv("TEST_GRPC"); !defined {
		t.SkipNow()
		return
	}
	suite.Run(t, new(NetworkControllerTestSuite))
}

func (suite *NetworkControllerTestSuite) TestNewNetworkController() {
	network := entity.Network{
		NodeName: suite.nodeName,
	}

	_, err := New(suite.kubectl, network)
	suite.NoError(err)
}

func (suite *NetworkControllerTestSuite) TestCreateNetwork() {
	//Parameters
	eth1 := entity.PhysicalPort{
		Name:     suite.ifName,
		MTU:      1500,
		VlanTags: []int{2043, 2143, 2243},
	}

	tName := namesgenerator.GetRandomName(0)
	network := entity.Network{
		BridgeName:    tName,
		BridgeType:    "ovs",
		NodeName:      suite.nodeName,
		PhysicalPorts: []entity.PhysicalPort{eth1},
	}

	nc, err := New(suite.kubectl, network)
	suite.NoError(err)
	err = nc.CreateNetwork()
	suite.NoError(err)

	defer exec.Command("ovs-vsctl", "del-br", tName).Run()
}
