package networkprovider

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	kc "github.com/linkernetworks/vortex/src/kubernetes"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

const DPDK_LOCAL_IP = "127.0.0.1"

func init() {
	rand.Seed(time.Now().UnixNano())
}

type DPDKNetworkTestSuite struct {
	suite.Suite
	sp             *serviceprovider.Container
	clusterNetwork entity.Network
	singleNetwork  entity.Network
}

func (suite *DPDKNetworkTestSuite) SetupSuite() {
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
					Address: DPDK_LOCAL_IP,
				},
			},
		},
	})
	suite.NoError(err)

	tName := namesgenerator.GetRandomName(0)
	suite.singleNetwork = entity.Network{
		Name: tName,
		OVSDPDK: entity.OVSDPDKNetwork{
			BridgeName:        tName,
			DPDKPhysicalPorts: []entity.DPDKPhysicalPort{},
		},
		Type:     entity.OVSUserspaceNetworkType,
		NodeName: nodeName,
	}

	suite.clusterNetwork = entity.Network{
		Name: tName,
		OVSDPDK: entity.OVSDPDKNetwork{
			BridgeName:        tName,
			DPDKPhysicalPorts: []entity.DPDKPhysicalPort{},
		},
		Type:        entity.OVSUserspaceNetworkType,
		NodeName:    nodeName,
		Clusterwise: true,
	}
}

func (suite *DPDKNetworkTestSuite) TearDownSuite() {
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
	suite.Run(t, new(DPDKNetworkTestSuite))
}

//Member funcion
func (suite *DPDKNetworkTestSuite) TestCreateOVSDPDKNetwork() {
	name := namesgenerator.GetRandomName(0)
	err := createOVSDPDKNetwork(DPDK_LOCAL_IP, name, []entity.DPDKPhysicalPort{})
	defer exec.Command("ovs-vsctl", "del-br", name).Run()
	suite.NoError(err)
}

func (suite *DPDKNetworkTestSuite) TestDeleteOVSDPDKNetwork() {
	name := namesgenerator.GetRandomName(0)
	// ovs-vsctl add-br br0 -- set bridge br0 datapath_type=netdev
	exec.Command("ovs-vsctl", "add-br", name, "--", "set", "bridge", name, "datapath_type=netdev").Run()
	err := deleteOVSDPDKNetwork(DPDK_LOCAL_IP, name)
	suite.NoError(err)
}

func (suite *DPDKNetworkTestSuite) TestCreateNetwork() {
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
			np = np.(OVSDPDKNetworkProvider)
			err = np.CreateNetwork(suite.sp, tc.network)
			suite.NoError(err)
			defer exec.Command("ovs-vsctl", "del-br", tc.network.OVS.BridgeName).Run()
		})
	}
}

func (suite *DPDKNetworkTestSuite) TestCreateNetworkFail() {
	network := entity.Network{
		Type: entity.OVSUserspaceNetworkType,
	}
	network.NodeName = "non-exist"
	np, err := GetNetworkProvider(&network)
	suite.NoError(err)
	np = np.(OVSDPDKNetworkProvider)
	err = np.CreateNetwork(suite.sp, &network)
	suite.Error(err)
}

func (suite *DPDKNetworkTestSuite) TestValidateBeforeCreating() {
	//Vlan
	//multiple network
	//single network

	//Prepare data
	eth1 := entity.DPDKPhysicalPort{
		Name:  namesgenerator.GetRandomName(0),
		MTU:   1500,
		PCIID: "0000:03:00.0",
		VlanTags: []int32{
			2043,
			2143,
			2243,
		},
	}

	tName := namesgenerator.GetRandomName(0)
	network := entity.Network{
		Name: tName,
		OVSDPDK: entity.OVSDPDKNetwork{
			BridgeName:        tName,
			DPDKPhysicalPorts: []entity.DPDKPhysicalPort{eth1},
		},
		Type: entity.OVSUserspaceNetworkType,
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
			np = np.(OVSDPDKNetworkProvider)

			err = np.ValidateBeforeCreating(suite.sp, tc.network)
			suite.NoError(err)
		})
	}
}

func (suite *DPDKNetworkTestSuite) TestValidateBeforeCreatingFail() {
	//Wrong Vlan
	//Wrong Case for multiple
	//Wrong Case for single

	//Prepare data
	eth1 := entity.DPDKPhysicalPort{
		Name:     namesgenerator.GetRandomName(0),
		MTU:      1500,
		PCIID:    "0000:03:00.0",
		VlanTags: []int32{2043, 2143, 22435},
	}

	tName := namesgenerator.GetRandomName(0)
	network := entity.Network{
		Name: tName,
		OVSDPDK: entity.OVSDPDKNetwork{
			BridgeName:        tName,
			DPDKPhysicalPorts: []entity.DPDKPhysicalPort{eth1},
		},
		Type: entity.OVSUserspaceNetworkType,
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
			np = np.(OVSDPDKNetworkProvider)

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

func (suite *DPDKNetworkTestSuite) TestDeleteNetwork() {
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
			np = np.(OVSDPDKNetworkProvider)
			err = np.CreateNetwork(suite.sp, tc.network)
			suite.NoError(err)

			// ovs-vsctl add-br br0 -- set bridge br0 datapath_type=netdev
			exec.Command("ovs-vsctl", "add-br", tc.network.OVSDPDK.BridgeName, "--", "set", "bridge", tc.network.OVSDPDK.BridgeName, "datapath_type=netdev").Run()
			//FIXME we need a function to check the bridge is exist
			err = np.DeleteNetwork(suite.sp, tc.network)
			suite.NoError(err)
		})
	}
}

func (suite *DPDKNetworkTestSuite) TestDeleteNetworkFail() {
	network := entity.Network{
		Type: entity.OVSUserspaceNetworkType,
	}
	network.NodeName = "non-exist"

	np, err := GetNetworkProvider(&network)
	suite.NoError(err)
	np = np.(OVSDPDKNetworkProvider)
	err = np.DeleteNetwork(suite.sp, &network)
	suite.Error(err)
}
