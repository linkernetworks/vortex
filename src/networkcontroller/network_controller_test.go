package networkcontroller

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/kubernetes"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeclientset "k8s.io/client-go/kubernetes/fake"

	"gopkg.in/mgo.v2/bson"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func execute(suite *suite.Suite, cmd *exec.Cmd) {
	w := bytes.NewBuffer(nil)
	cmd.Stderr = w
	err := cmd.Run()
	suite.NoError(err)
	fmt.Printf("Stderr: %s\n", string(w.Bytes()))
}

type NetworkControllerTestSuite struct {
	suite.Suite
	kubectl  *kubernetes.KubeCtl
	ifName   string
	nodeName string
}

func (suite *NetworkControllerTestSuite) SetupSuite() {
	// init fakeclient
	fakeclient := fakeclientset.NewSimpleClientset()
	suite.kubectl = kubernetes.New(fakeclient)

	//Create a fake clinet
	//Init
	nodeAddr := corev1.NodeAddress{
		Type:    "InternalIP",
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
	suite.ifName = bson.NewObjectId().Hex()[12:24]
	pName := bson.NewObjectId().Hex()[12:24]
	//Create a veth for testing
	fmt.Println("ip", "link", "add", suite.ifName, "type", "veth", "peer", "name", pName)
	cmd := exec.Command("ip", "link", "add", suite.ifName, "type", "veth", "peer", "name", pName)
	execute(&suite.Suite, cmd)
}

func (suite *NetworkControllerTestSuite) TearDownSuite() {
	cmd := exec.Command("ip", "link", "del", suite.ifName)
	execute(&suite.Suite, cmd)
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

func (suite *NetworkControllerTestSuite) TestNew() {
	_, err := New(net.JoinHostPort("127.0.0.1", DEFAULT_CONTROLLER_PORT))
	suite.NoError(err)
}

func (suite *NetworkControllerTestSuite) TestCreateNetwork() {
	tName := namesgenerator.GetRandomName(0)
	network := entity.Network{
		Type:       entity.OVSKernelspaceNetworkType,
		Name:       tName,
		VLANTags:   []int32{0, 2048, 4095},
		BridgeName: tName,
		Nodes: []entity.Node{
			entity.Node{
				Name: suite.nodeName,
				PhyInterfaces: []entity.PhyInterface{
					entity.PhyInterface{
						Name:  suite.ifName,
						PCIID: "",
					},
				},
			},
		},
	}

	nodeIP, err := suite.kubectl.GetNodeInternalIP(suite.nodeName)
	suite.NoError(err)
	nc, err := New(net.JoinHostPort(nodeIP, DEFAULT_CONTROLLER_PORT))
	suite.NoError(err)
	err = nc.CreateOVSNetwork("system", tName, network.Nodes[0].PhyInterfaces, network.VLANTags)
	suite.NoError(err)

	//TODO we need support the list function to check the ovs is existed
	defer exec.Command("ovs-vsctl", "del-br", tName).Run()
}

func (suite *NetworkControllerTestSuite) TestCreateOVSUserpsaceNetwork() {
	tName := namesgenerator.GetRandomName(0)
	network := entity.Network{
		Type:       entity.OVSUserspaceNetworkType,
		IsDPDKPort: false,
		Name:       tName,
		VLANTags:   []int32{0, 2048, 4095},
		BridgeName: tName,
		Nodes: []entity.Node{
			entity.Node{
				Name: suite.nodeName,
				PhyInterfaces: []entity.PhyInterface{
					entity.PhyInterface{
						Name:  suite.ifName,
						PCIID: "",
					},
				},
			},
		},
	}

	nodeIP, err := suite.kubectl.GetNodeInternalIP(suite.nodeName)
	suite.NoError(err)
	nc, err := New(net.JoinHostPort(nodeIP, DEFAULT_CONTROLLER_PORT))
	suite.NoError(err)
	err = nc.CreateOVSNetwork("netdev", tName, network.Nodes[0].PhyInterfaces, network.VLANTags)
	suite.NoError(err)

	//TODO we need support the list function to check the ovs is existed
	defer exec.Command("ovs-vsctl", "del-br", tName).Run()
}

func (suite *NetworkControllerTestSuite) TestCreateOVSDPDKNetwork() {
	tName := namesgenerator.GetRandomName(0)
	network := entity.Network{
		Type:       entity.OVSUserspaceNetworkType,
		IsDPDKPort: true,
		Name:       tName,
		VLANTags:   []int32{0, 2048, 4095},
		BridgeName: tName,
		Nodes: []entity.Node{
			entity.Node{
				Name: suite.nodeName,
				PhyInterfaces: []entity.PhyInterface{
					entity.PhyInterface{
						Name:  suite.ifName,
						PCIID: "0000:00:08.0",
					},
				},
			},
		},
	}

	nodeIP, err := suite.kubectl.GetNodeInternalIP(suite.nodeName)
	suite.NoError(err)
	nc, err := New(net.JoinHostPort(nodeIP, DEFAULT_CONTROLLER_PORT))
	suite.NoError(err)
	err = nc.CreateOVSNetwork("netdev", tName, network.Nodes[0].PhyInterfaces, network.VLANTags)
	suite.NoError(err)

	//TODO we need support the list function to check the ovs is existed
	defer exec.Command("ovs-vsctl", "del-br", tName).Run()
}

func (suite *NetworkControllerTestSuite) TestDeleteNetwork() {
	tName := namesgenerator.GetRandomName(0)
	network := entity.Network{
		Type:       entity.OVSKernelspaceNetworkType,
		Name:       tName,
		VLANTags:   []int32{0, 2048, 4095},
		BridgeName: tName,
		Nodes: []entity.Node{
			entity.Node{
				Name: suite.nodeName,
				PhyInterfaces: []entity.PhyInterface{
					entity.PhyInterface{
						Name:  suite.ifName,
						PCIID: "",
					},
				},
			},
		},
	}

	nodeIP, err := suite.kubectl.GetNodeInternalIP(suite.nodeName)
	suite.NoError(err)
	nc, err := New(net.JoinHostPort(nodeIP, DEFAULT_CONTROLLER_PORT))
	suite.NoError(err)
	err = nc.CreateOVSNetwork("system", tName, network.Nodes[0].PhyInterfaces, network.VLANTags)
	suite.NoError(err)

	err = nc.DeleteOVSNetwork(tName)
	suite.NoError(err)
}

func (suite *NetworkControllerTestSuite) TestCreateNetworkWithInvalidAddress() {
	nc, err := New(net.JoinHostPort("a.b.c.d", DEFAULT_CONTROLLER_PORT))
	suite.NoError(err)

	tName := namesgenerator.GetRandomName(0)
	err = nc.CreateOVSNetwork("system", tName, []entity.PhyInterface{}, []int32{})
	suite.Error(err)
}
