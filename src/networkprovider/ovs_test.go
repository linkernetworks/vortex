package networkprovider

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	//"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	kc "github.com/linkernetworks/vortex/src/kubernetes"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"

	//mgo "gopkg.in/mgo.v2"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

const LOCAL_IP = "127.0.0.1"

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

type NetworkTestSuite struct {
	suite.Suite
	sp      *serviceprovider.Container
	np      OVSNetworkProvider
	network entity.Network
}

func (suite *NetworkTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	suite.sp = serviceprovider.NewForTesting(cf)

	//init fakeclient
	fakeclient := fakeclientset.NewSimpleClientset()
	namespace := "default"
	suite.sp.KubeCtl = kc.New(fakeclient, namespace)

	//Create a fake clinet
	//Init node
	nodeName := namesgenerator.GetRandomName(0)
	_, err := suite.sp.KubeCtl.Clientset.CoreV1().Nodes().Create(&corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: nodeName,
		},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{
				{
					Type:    "ExternalIP",
					Address: LOCAL_IP,
				},
			},
		},
	})
	suite.NoError(err)

	tName := namesgenerator.GetRandomName(0)
	suite.network = entity.Network{
		Name: tName,
		OVS: entity.OVSNetwork{
			BridgeName:    tName,
			PhysicalPorts: []entity.PhysicalPort{},
		},
		Type:     "ovs",
		NodeName: nodeName,
	}

	np, err := GetNetworkProvider(&suite.network)
	suite.NoError(err)
	suite.np = np.(OVSNetworkProvider)
}

func (suite *NetworkTestSuite) TearDownSuite() {
}

func TestNetworkSuite(t *testing.T) {
	if runtime.GOOS != "linux" {
		fmt.Println("We only testing the ovs function on Linux Host")
		t.Skip()
		return
	}
	if _, defined := os.LookupEnv("TEST_GRPC"); !defined {
		t.SkipNow()
		return
	}
	suite.Run(t, new(NetworkTestSuite))
}

func (suite *NetworkTestSuite) TestCreateOVSNetwork() {
	err := createOVSNetwork(LOCAL_IP, suite.np.BridgeName, []entity.PhysicalPort{})
	suite.NoError(err)
	defer exec.Command("ovs-vsctl", "del-br", suite.np.BridgeName).Run()
}

func (suite *NetworkTestSuite) TestCreateNetwork() {
	//Parameters
	err := suite.np.CreateNetwork(suite.sp, suite.network)
	suite.NoError(err)
	defer exec.Command("ovs-vsctl", "del-br", suite.np.BridgeName).Run()
}

func (suite *NetworkTestSuite) TestCreateNetworkWithCluster() {
	//Parameters
	network := entity.Network{
		OVS: entity.OVSNetwork{
			BridgeName: suite.np.BridgeName,
		},
		Clusterwise: true,
	}
	err := suite.np.CreateNetwork(suite.sp, network)
	suite.NoError(err)
	defer exec.Command("ovs-vsctl", "del-br", suite.np.BridgeName).Run()
}

func (suite *NetworkTestSuite) TestCreateNetworkFail() {
	network := entity.Network{}
	network.NodeName = "non-exist"
	err := suite.np.CreateNetwork(suite.sp, network)
	suite.Error(err)
}

func (suite *NetworkTestSuite) TestValidateBeforeCreating() {
	//Parameters
	eth1 := entity.PhysicalPort{
		Name:     namesgenerator.GetRandomName(0),
		MTU:      1500,
		VlanTags: []int{2043, 2143, 2243},
	}

	ovsProvider := suite.np
	ovsProvider.PhysicalPorts = []entity.PhysicalPort{eth1}
	err := ovsProvider.ValidateBeforeCreating(suite.sp, suite.network)
	suite.NoError(err)
}

func (suite *NetworkTestSuite) TestValidateBeforeCreatingFail() {
	//Parameters
	ovsProvider := suite.np

	//create a mongo-document to test duplicated name
	session := suite.sp.Mongo.NewSession()
	err := session.C(entity.NetworkCollectionName).Insert(suite.network)
	defer session.C(entity.NetworkCollectionName).Remove(suite.network)
	suite.NoError(err)
	err = ovsProvider.ValidateBeforeCreating(suite.sp, suite.network)
	suite.Error(err)

	//Test wrong vlan ID
	eth1 := entity.PhysicalPort{
		Name:     namesgenerator.GetRandomName(0),
		MTU:      1500,
		VlanTags: []int{2043, 2143, 22434},
	}

	ovsProvider.PhysicalPorts = []entity.PhysicalPort{eth1}
	err = ovsProvider.ValidateBeforeCreating(suite.sp, suite.network)
	suite.Error(err)
}

func (suite *NetworkTestSuite) TestDeleteNetwork() {
	//Parameters
	exec.Command("ovs-vsctl", "add-br", suite.np.BridgeName).Run()
	//FIXME we need a function to check the bridge is exist
	err := suite.np.DeleteNetwork(suite.sp, suite.network)
	suite.NoError(err)
}

func (suite *NetworkTestSuite) TestDeleteNetworkFail() {
	network := entity.Network{}
	network.NodeName = "non-exist"
	err := suite.np.DeleteNetwork(suite.sp, network)
	suite.Error(err)
}
