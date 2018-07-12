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

const OVS_LOCAL_IP = "127.0.0.1"

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

//OK
type OVSSystemNetworkTestSuite struct {
	suite.Suite
	sp                 *serviceprovider.Container
	standaloneNetwork  entity.Network
	clusterwiseNetwork entity.Network
}

// OK
func (suite *OVSSystemNetworkTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	suite.sp = serviceprovider.NewForTesting(cf)

	// init fakeclient
	fakeclient := fakeclientset.NewSimpleClientset()
	namespace := "default"
	suite.sp.KubeCtl = kc.New(fakeclient, namespace)

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
					Type:    "ExternalIP",
					Address: OVS_LOCAL_IP,
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
					Type:    "ExternalIP",
					Address: OVS_LOCAL_IP,
				},
			},
		},
	})
	suite.NoError(err)

	tName := namesgenerator.GetRandomName(0)

	suite.standaloneNetwork = entity.Network{
		Type:       entity.OVSKernelspaceNetworkType,
		Name:       tName,
		BridgeName: tName,
		Nodes: []entity.Node{
			entity.Node{
				Name: nodeName1,
			},
		},
	}

	suite.clusterwiseNetwork = entity.Network{
		Type:       entity.OVSKernelspaceNetworkType,
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

func (suite *OVSSystemNetworkTestSuite) TearDownSuite() {}

// OK
func TestOVSNetworkSuite(t *testing.T) {
	if runtime.GOOS != "linux" {
		fmt.Println("We only testing the ovs function on Linux Host")
		t.Skip()
		return
	}
	if _, defined := os.LookupEnv("TEST_GRPC"); !defined {
		t.SkipNow()
		return
	}
	suite.Run(t, new(OVSSystemNetworkTestSuite))
}

// OK
func (suite *OVSSystemNetworkTestSuite) TestCreateOVSNetwork() {
	brName := namesgenerator.GetRandomName(0)
	err := createOVSNetwork(
		OVS_LOCAL_IP,
		brName,
		[]entity.PhyInterface{},
		[]int32{0, 2048, 4095},
	)
	defer exec.Command("ovs-vsctl", "del-br", brName).Run()
	suite.NoError(err)
}

// OK
func (suite *OVSSystemNetworkTestSuite) TestDeleteOVSNetwork() {
	brName := namesgenerator.GetRandomName(0)
	// ovs-vsctl add-br br0 -- set bridge br0 datapath_type=netdev
	exec.Command("ovs-vsctl", "add-br", brName, "--", "set", "bridge", brName, "datapath_type=netdev").Run()
	err := deleteOVSNetwork(OVS_LOCAL_IP, brName)
	suite.NoError(err)
}

func (suite *OVSSystemNetworkTestSuite) TestCreateNetwork() {
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
			np = np.(kernelspaceNetworkProvider)
			err = np.CreateNetwork(suite.sp)
			suite.NoError(err)
			defer exec.Command("ovs-vsctl", "del-br", tc.network.BridgeName).Run()
		})
	}
}

func (suite *OVSSystemNetworkTestSuite) TestCreateNetworkFail() {
	network := entity.Network{
		Type:       entity.OVSKernelspaceNetworkType,
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
	np = np.(kernelspaceNetworkProvider)
	err = np.CreateNetwork(suite.sp)
	suite.Error(err)
}

func (suite *OVSSystemNetworkTestSuite) TestValidateBeforeCreating() {
	tName := namesgenerator.GetRandomName(0)
	// Valid VLAN number
	validVLANnetwork := entity.Network{
		Type:       entity.OVSKernelspaceNetworkType,
		Name:       tName,
		VLANTags:   []int32{0, 2048, 4095},
		BridgeName: "bro",
		Nodes: []entity.Node{
			entity.Node{
				Name: "node-1",
				PhyInterfaces: []entity.PhyInterface{
					entity.PhyInterface{
						Name:  namesgenerator.GetRandomName(0),
						PCIID: "",
					},
				},
			},
		},
	}

	testCases := []struct {
		caseName string
		network  *entity.Network
	}{
		{"validVLAN", &validVLANnetwork},
		{"singelNetwork", &suite.standaloneNetwork},
		{"clusterNetwork", &suite.clusterwiseNetwork},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.caseName, func(t *testing.T) {
			//Parameters
			np, err := GetNetworkProvider(tc.network)
			suite.NoError(err)
			np = np.(kernelspaceNetworkProvider)

			err = np.ValidateBeforeCreating(suite.sp)
			suite.NoError(err)
		})
	}
}

// OK
func (suite *OVSSystemNetworkTestSuite) TestValidateBeforeCreatingFail() {
	tName := namesgenerator.GetRandomName(0)
	// Invalid VLAN number
	invalidVLANnetwork := entity.Network{
		Type:       entity.OVSKernelspaceNetworkType,
		Name:       tName,
		VLANTags:   []int32{22435, 2143, 4096},
		BridgeName: "brx",
		Nodes: []entity.Node{
			entity.Node{
				Name: "node-1",
				PhyInterfaces: []entity.PhyInterface{
					entity.PhyInterface{
						Name:  namesgenerator.GetRandomName(0),
						PCIID: "",
					},
				},
			},
		},
	}

	// Invalid Case for duplicated name in standalone network
	// Invalid Case for duplicates name in clusterwise network
	testCases := []struct {
		caseName string
		network  *entity.Network
		mongo    bool
	}{
		{"invalidVLAN", &invalidVLANnetwork, false},
		{"singelNetwork", &suite.standaloneNetwork, true},
		{"clusterNetwork", &suite.clusterwiseNetwork, true},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.caseName, func(t *testing.T) {
			//Parameters
			np, err := GetNetworkProvider(tc.network)
			suite.NoError(err)
			np = np.(kernelspaceNetworkProvider)

			if tc.mongo {
				//create a mongo-document to test duplicated name
				session := suite.sp.Mongo.NewSession()
				err := session.C(entity.NetworkCollectionName).Insert(tc.network)
				defer session.C(entity.NetworkCollectionName).Remove(tc.network)
				suite.NoError(err)
			}
			err = np.ValidateBeforeCreating(suite.sp)
			suite.Error(err)
		})
	}
}

func (suite *OVSSystemNetworkTestSuite) TestDeleteNetwork() {
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
			np = np.(kernelspaceNetworkProvider)
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

// OK
func (suite *OVSSystemNetworkTestSuite) TestDeleteNetworkFail() {
	network := entity.Network{
		Type:       entity.OVSKernelspaceNetworkType,
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
	np = np.(kernelspaceNetworkProvider)
	err = np.DeleteNetwork(suite.sp)
	suite.Error(err)
}
