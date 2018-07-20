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

type OVSNetdevNetworkTestSuite struct {
	suite.Suite
	sp                 *serviceprovider.Container
	standaloneNetwork  entity.Network
	clusterwiseNetwork entity.Network
}

func (suite *OVSNetdevNetworkTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	suite.sp = serviceprovider.NewForTesting(cf)

	// init fakeclient
	fakeclient := fakeclientset.NewSimpleClientset()
	suite.sp.KubeCtl = kc.New(fakeclient)

	// Create a fake clinet
	// Initial nodes
	nodeName1 := namesgenerator.GetRandomName(0)
	nodeName2 := namesgenerator.GetRandomName(1)
	_, err := suite.sp.KubeCtl.Clientset.CoreV1().Nodes().Create(&corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: nodeName1,
		},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{
				{
					Type:    "InternalIP",
					Address: DPDK_LOCAL_IP,
				},
			},
		},
	})
	suite.NoError(err)

	_, err = suite.sp.KubeCtl.Clientset.CoreV1().Nodes().Create(&corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: nodeName2,
		},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{
				{
					Type:    "InternalIP",
					Address: DPDK_LOCAL_IP,
				},
			},
		},
	})
	suite.NoError(err)

	tName := namesgenerator.GetRandomName(0)

	suite.standaloneNetwork = entity.Network{
		Type:       entity.OVSUserspaceNetworkType,
		IsDPDKPort: false,
		Name:       tName,
		BridgeName: tName,
		Nodes: []entity.Node{
			entity.Node{
				Name: nodeName1,
			},
		},
	}

	suite.clusterwiseNetwork = entity.Network{
		Type:       entity.OVSUserspaceNetworkType,
		IsDPDKPort: false,
		Name:       tName,
		BridgeName: tName,
		Nodes: []entity.Node{
			entity.Node{
				Name: nodeName1,
			},
			entity.Node{
				Name: nodeName2,
			},
		},
	}
}

func (suite *OVSNetdevNetworkTestSuite) TearDownSuite() {}

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
	suite.Run(t, new(OVSNetdevNetworkTestSuite))
}

func (suite *OVSNetdevNetworkTestSuite) TestCreateOVSDPDKNetwork() {
	brName := namesgenerator.GetRandomName(0)
	err := createOVSDPDKNetwork(
		DPDK_LOCAL_IP,
		brName,
		[]entity.PhyInterface{},
		[]int32{0, 2048, 4095},
	)
	defer exec.Command("ovs-vsctl", "del-br", brName).Run()
	suite.NoError(err)
}

func (suite *OVSNetdevNetworkTestSuite) TestCreateOVSUserspaceNetwork() {
	brName := namesgenerator.GetRandomName(0)
	err := createOVSUserspaceNetwork(
		DPDK_LOCAL_IP,
		brName,
		[]entity.PhyInterface{},
		[]int32{0, 2048, 4095},
	)
	defer exec.Command("ovs-vsctl", "del-br", brName).Run()
	suite.NoError(err)
}

func (suite *OVSNetdevNetworkTestSuite) TestDeleteOVSUserspaceNetwork() {
	brName := namesgenerator.GetRandomName(0)
	// ovs-vsctl add-br br0 -- set bridge br0 datapath_type=netdev
	exec.Command("ovs-vsctl", "add-br", brName, "--", "set", "bridge", brName, "datapath_type=netdev").Run()
	err := deleteOVSUserspaceNetwork(DPDK_LOCAL_IP, brName)
	suite.NoError(err)
}

func (suite *OVSNetdevNetworkTestSuite) TestCreateNetwork() {
	testCases := []struct {
		caseName string
		network  *entity.Network
	}{
		{"standaloneNetwork", &suite.standaloneNetwork},
		{"clusterwiseNetwork", &suite.clusterwiseNetwork},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.caseName, func(t *testing.T) {
			np, err := GetNetworkProvider(tc.network)
			suite.NoError(err)
			np = np.(userspaceNetworkProvider)
			err = np.CreateNetwork(suite.sp)
			suite.NoError(err)
			defer exec.Command("ovs-vsctl", "del-br", tc.network.BridgeName).Run()
		})
	}
}

func (suite *OVSNetdevNetworkTestSuite) TestCreateNetworkFail() {
	network := entity.Network{
		Type:       entity.OVSUserspaceNetworkType,
		IsDPDKPort: false,
		Name:       "none-exist-network",
		BridgeName: "none",
		Nodes: []entity.Node{
			entity.Node{
				Name: namesgenerator.GetRandomName(0),
			},
		},
	}
	np, err := GetNetworkProvider(&network)
	suite.NoError(err)
	np = np.(userspaceNetworkProvider)
	err = np.CreateNetwork(suite.sp)
	suite.Error(err)
}

func (suite *OVSNetdevNetworkTestSuite) TestDeleteNetwork() {
	testCases := []struct {
		caseName string
		network  *entity.Network
	}{
		{"standaloneNetwork", &suite.standaloneNetwork},
		{"clusterwiseNetwork", &suite.clusterwiseNetwork},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.caseName, func(t *testing.T) {
			//Parameters
			np, err := GetNetworkProvider(tc.network)
			suite.NoError(err)
			np = np.(userspaceNetworkProvider)
			err = np.CreateNetwork(suite.sp)
			suite.NoError(err)

			// ovs-vsctl add-br br0 -- set bridge br0 datapath_type=netdev
			exec.Command("ovs-vsctl", "add-br", tc.network.BridgeName, "--", "set", "bridge", tc.network.BridgeName, "datapath_type=netdev").Run()
			//FIXME we need a function to check the bridge is exist
			err = np.DeleteNetwork(suite.sp)
			suite.NoError(err)
		})
	}
}

func (suite *OVSNetdevNetworkTestSuite) TestDeleteNetworkFail() {
	network := entity.Network{
		Type:       entity.OVSUserspaceNetworkType,
		IsDPDKPort: false,
		Name:       "none-exist-network",
		BridgeName: "none",
		Nodes: []entity.Node{
			entity.Node{
				Name: namesgenerator.GetRandomName(0),
			},
		},
	}
	np, err := GetNetworkProvider(&network)
	suite.NoError(err)
	np = np.(userspaceNetworkProvider)
	err = np.DeleteNetwork(suite.sp)
	suite.Error(err)
}
