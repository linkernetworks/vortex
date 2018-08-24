package ovscontroller

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
	_ "github.com/linkernetworks/vortex/src/entity"
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

type OVSControllerTestSuite struct {
	suite.Suite
	sp         *serviceprovider.Container
	nodeName   string
	bridgeName string
}

func (suite *OVSControllerTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	suite.sp = serviceprovider.NewForTesting(cf)

	// init fakeclient
	fakeclient := fakeclientset.NewSimpleClientset()
	suite.sp.KubeCtl = kc.New(fakeclient)

	suite.bridgeName = namesgenerator.GetRandomName(0)[0:6]

	// Create a fake clinet
	// Initial nodes
	suite.nodeName = namesgenerator.GetRandomName(0)
	_, err := suite.sp.KubeCtl.Clientset.CoreV1().Nodes().Create(&corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: suite.nodeName,
		},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{
				{
					Type:    "InternalIP",
					Address: OVS_LOCAL_IP,
				},
			},
		},
	})
	suite.NoError(err)

	execute(&suite.Suite, exec.Command("ovs-vsctl", "add-br", suite.bridgeName))
}

func (suite *OVSControllerTestSuite) TearDownSuite() {
	defer exec.Command("ovs-vsctl", "del-br", suite.bridgeName).Run()
}

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
	suite.Run(t, new(OVSControllerTestSuite))
}

// OK
func (suite *OVSControllerTestSuite) TestDumpOVSPorts() {
	portStats, err := DumpPorts(suite.sp, suite.nodeName, suite.bridgeName)
	suite.NoError(err)
	suite.Equal(1, len(portStats))
}
