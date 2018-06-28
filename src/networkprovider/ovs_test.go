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
	//"gopkg.in/mgo.v2/bson"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
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

type NetworkTestSuite struct {
	suite.Suite
	sp       *serviceprovider.Container
	fakeName string //Use for non-connectivity node
	np       OVSNetworkProvider
	network  *entity.Network
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
					Address: "127.0.0.1",
				},
			},
		},
	})
	suite.NoError(err)

	suite.fakeName = namesgenerator.GetRandomName(0)
	_, err = suite.sp.KubeCtl.Clientset.CoreV1().Nodes().Create(&corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: suite.fakeName,
		},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{
				{
					Type:    "ExternalIP",
					Address: "1.2.3.4",
				},
			},
		},
	})
	suite.NoError(err)

	tName := namesgenerator.GetRandomName(0)
	suite.network = &entity.Network{
		Name: tName,
		OVS: entity.OVSNetwork{
			BridgeName:    tName,
			PhysicalPorts: []entity.PhysicalPort{},
		},
		Type:     "ovs",
		NodeName: nodeName,
	}

	np, err := GetNetworkProvider(suite.network)
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

func (suite *NetworkTestSuite) TestCreateNetwork() {
	//Parameters
	err := suite.np.CreateNetwork(suite.sp, *(suite.network))
	suite.NoError(err)
	defer exec.Command("ovs-vsctl", "del-br", suite.np.BridgeName).Run()
}

func (suite *NetworkTestSuite) TestCreateNetworkFail() {
	network := suite.network
	network.NodeName = "non-exist"
	err := suite.np.CreateNetwork(suite.sp, *network)
	suite.Error(err)

	network.NodeName = suite.fakeName
	err = suite.np.CreateNetwork(suite.sp, *network)
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
	err := ovsProvider.ValidateBeforeCreating(suite.sp, *suite.network)
	suite.NoError(err)
}

func (suite *NetworkTestSuite) TestValidateBeforeCreatingFail() {
	//Parameters
	eth1 := entity.PhysicalPort{
		Name:     namesgenerator.GetRandomName(0),
		MTU:      1500,
		VlanTags: []int{2043, 2143, 22434},
	}

	ovsProvider := suite.np
	ovsProvider.PhysicalPorts = []entity.PhysicalPort{eth1}
	err := ovsProvider.ValidateBeforeCreating(suite.sp, *suite.network)
	suite.Error(err)
}
