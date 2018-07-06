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
	sp             *serviceprovider.Container
	clusterNetwork entity.Network
	singleNetwork  entity.Network
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
	suite.singleNetwork = entity.Network{
		Name: tName,
		OVS: entity.OVSNetwork{
			BridgeName:    tName,
			PhysicalPorts: []entity.PhysicalPort{},
		},
		Type:     "ovs",
		NodeName: nodeName,
	}

	suite.clusterNetwork = entity.Network{
		Name: tName,
		OVS: entity.OVSNetwork{
			BridgeName:    tName,
			PhysicalPorts: []entity.PhysicalPort{},
		},
		Type:        "ovs",
		NodeName:    nodeName,
		Clusterwise: true,
	}

	//np, err := GetNetworkProvider(&suite.network)
	//suite.NoError(err)
	//suite.np = np.(OVSNetworkProvider)
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

//Member funcion
func (suite *NetworkTestSuite) TestCreateOVSNetwork() {
	name := namesgenerator.GetRandomName(0)
	err := createOVSNetwork(LOCAL_IP, name, []entity.PhysicalPort{})
	defer exec.Command("ovs-vsctl", "del-br", name).Run()
	suite.NoError(err)
}

func (suite *NetworkTestSuite) TestDeleteOVSNetwork() {
	name := namesgenerator.GetRandomName(0)
	exec.Command("ovs-vsctl", "add-br", name).Run()
	err := deleteOVSNetwork(LOCAL_IP, name)
	suite.NoError(err)
}

func (suite *NetworkTestSuite) TestCreateNetwork() {
	testCases := []struct {
		caseName string
		network  *entity.Network
	}{
		{"singelNetwork", &suite.singleNetwork},
		{"clusterNetwork", &suite.clusterNetwork},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.caseName, func(t *testing.T) {
			//Parameters
			np, err := GetNetworkProvider(tc.network)
			suite.NoError(err)
			np = np.(OVSNetworkProvider)
			err = np.CreateNetwork(suite.sp, tc.network)
			suite.NoError(err)
			defer exec.Command("ovs-vsctl", "del-br", tc.network.OVS.BridgeName).Run()
		})
	}
}

func (suite *NetworkTestSuite) TestCreateNetworkFail() {
	network := entity.Network{
		Type: "ovs",
	}
	network.NodeName = "non-exist"
	np, err := GetNetworkProvider(&network)
	suite.NoError(err)
	np = np.(OVSNetworkProvider)
	err = np.CreateNetwork(suite.sp, &network)
	suite.Error(err)
}

func (suite *NetworkTestSuite) TestValidateBeforeCreating() {
	//Vlan
	//multiple network
	//single network

	//Prepare data
	eth1 := entity.PhysicalPort{
		Name:     namesgenerator.GetRandomName(0),
		MTU:      1500,
		VlanTags: []int32{2043, 2143, 2243},
	}

	tName := namesgenerator.GetRandomName(0)
	network := entity.Network{
		Name: tName,
		OVS: entity.OVSNetwork{
			BridgeName:    tName,
			PhysicalPorts: []entity.PhysicalPort{eth1},
		},
		Type: "ovs",
	}

	testCases := []struct {
		caseName string
		network  *entity.Network
	}{
		{"valid", &network},
		{"singelNetwork", &suite.singleNetwork},
		{"clusterNetwork", &suite.clusterNetwork},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.caseName, func(t *testing.T) {
			//Parameters
			np, err := GetNetworkProvider(tc.network)
			suite.NoError(err)
			np = np.(OVSNetworkProvider)

			err = np.ValidateBeforeCreating(suite.sp, tc.network)
			suite.NoError(err)
		})
	}
}

func (suite *NetworkTestSuite) TestValidateBeforeCreatingFail() {
	//Wrong Vlan
	//Wrong Case for multiple
	//Wrong Case for single

	//Prepare data
	eth1 := entity.PhysicalPort{
		Name:     namesgenerator.GetRandomName(0),
		MTU:      1500,
		VlanTags: []int32{2043, 2143, 22435},
	}

	tName := namesgenerator.GetRandomName(0)
	network := entity.Network{
		Name: tName,
		OVS: entity.OVSNetwork{
			BridgeName:    tName,
			PhysicalPorts: []entity.PhysicalPort{eth1},
		},
		Type: "ovs",
	}

	testCases := []struct {
		caseName string
		network  *entity.Network
		mongo    bool
	}{
		{"invalidValid", &network, false},
		{"singelNetwork", &suite.singleNetwork, true},
		{"clusterNetwork", &suite.clusterNetwork, true},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.caseName, func(t *testing.T) {
			//Parameters
			np, err := GetNetworkProvider(tc.network)
			suite.NoError(err)
			np = np.(OVSNetworkProvider)

			if tc.mongo {
				//create a mongo-document to test duplicated name
				session := suite.sp.Mongo.NewSession()
				err := session.C(entity.NetworkCollectionName).Insert(tc.network)
				defer session.C(entity.NetworkCollectionName).Remove(tc.network)
				suite.NoError(err)
			}
			err = np.ValidateBeforeCreating(suite.sp, tc.network)
			suite.Error(err)
		})
	}
}

func (suite *NetworkTestSuite) TestDeleteNetwork() {
	testCases := []struct {
		caseName string
		network  *entity.Network
	}{
		{"singelNetwork", &suite.singleNetwork},
		{"clusterNetwork", &suite.clusterNetwork},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.caseName, func(t *testing.T) {
			//Parameters
			np, err := GetNetworkProvider(tc.network)
			suite.NoError(err)
			np = np.(OVSNetworkProvider)
			err = np.CreateNetwork(suite.sp, tc.network)
			suite.NoError(err)

			exec.Command("ovs-vsctl", "add-br", tc.network.OVS.BridgeName).Run()
			//FIXME we need a function to check the bridge is exist
			err = np.DeleteNetwork(suite.sp, tc.network)
			suite.NoError(err)
		})
	}
}

func (suite *NetworkTestSuite) TestDeleteNetworkFail() {
	network := entity.Network{
		Type: "ovs",
	}
	network.NodeName = "non-exist"

	np, err := GetNetworkProvider(&network)
	suite.NoError(err)
	np = np.(OVSNetworkProvider)
	err = np.DeleteNetwork(suite.sp, &network)
	suite.Error(err)
}
