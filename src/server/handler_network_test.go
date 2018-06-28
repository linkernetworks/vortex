package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"

	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	kc "github.com/linkernetworks/vortex/src/kubernetes"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

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
	wc       *restful.Container
	session  *mongo.Session
	ifName   string
	nodeName string
}

func (suite *NetworkTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	//init session
	suite.session = sp.Mongo.NewSession()
	//init restful container
	suite.wc = restful.NewContainer()
	service := newNetworkService(sp)
	suite.wc.Add(service)

	//init fakeclient
	fakeclient := fakeclientset.NewSimpleClientset()
	namespace := "default"
	sp.KubeCtl = kc.New(fakeclient, namespace)

	//Create a fake clinet
	//Init
	nodeAddr := corev1.NodeAddress{
		Type:    "ExternalIP",
		Address: "127.0.0.1",
	}

	suite.nodeName = namesgenerator.GetRandomName(0)
	node := corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: suite.nodeName,
		},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{nodeAddr},
		},
	}
	_, err := sp.KubeCtl.Clientset.CoreV1().Nodes().Create(&node)
	suite.NoError(err)

	//There's a length limit of link name
	suite.ifName = bson.NewObjectId().Hex()[12:24]
	pName := bson.NewObjectId().Hex()[12:24]
	//Create a veth for testing
	fmt.Println("ip", "link", "add", suite.ifName, "type", "veth", "peer", "name", pName)
	cmd := exec.Command("ip", "link", "add", suite.ifName, "type", "veth", "peer", "name", pName)
	execute(&suite.Suite, cmd)
}

func (suite *NetworkTestSuite) TearDownSuite() {
	cmd := exec.Command("ip", "link", "del", suite.ifName)
	execute(&suite.Suite, cmd)
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

	bodyBytes, err := json.MarshalIndent(network, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/networks", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	//assertResponseCode(t, http.StatusOK, httpWriter)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
	defer suite.session.Remove(entity.NetworkCollectionName, "bridgeName", tName)
	defer exec.Command("ovs-vsctl", "del-br", tName).Run()

	//We use the new write but empty input
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
	//Create again and it should fail since the name exist
	bodyReader = strings.NewReader(string(bodyBytes))
	httpRequest, err = http.NewRequest("POST", "http://localhost:7890/v1/networks", bodyReader)
	suite.NoError(err)
	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusConflict, httpWriter)
}

func (suite *NetworkTestSuite) TestWrongVlangTag() {
	tName := namesgenerator.GetRandomName(0)
	network := entity.Network{
		BridgeName: tName,
		BridgeType: "ovs",
		NodeName:   suite.nodeName,
		PhysicalPorts: []entity.PhysicalPort{
			{
				Name:     suite.ifName,
				MTU:      1500,
				VlanTags: []int{1234, 2143, 50000},
			},
		}}

	bodyBytes, err := json.MarshalIndent(network, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/networks", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
}

func (suite *NetworkTestSuite) TestDeleteNetwork() {
	tName := namesgenerator.GetRandomName(0)
	network := entity.Network{
		ID:            bson.NewObjectId(),
		BridgeName:    tName,
		BridgeType:    "ovs",
		NodeName:      suite.nodeName,
		PhysicalPorts: []entity.PhysicalPort{},
	}

	//Create data into mongo manually
	suite.session.C(entity.NetworkCollectionName).Insert(network)
	err := exec.Command("ovs-vsctl", "add-br", tName).Run()
	suite.NoError(err)
	defer suite.session.Remove(entity.NetworkCollectionName, "bridgeName", tName)

	httpRequestDelete, err := http.NewRequest("DELETE", "http://localhost:7890/v1/networks/"+network.ID.Hex(), nil)
	httpWriterDelete := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriterDelete, httpRequestDelete)
	assertResponseCode(suite.T(), http.StatusOK, httpWriterDelete)
	err = suite.session.FindOne(entity.NetworkCollectionName, bson.M{"_id": network.ID}, &network)
	suite.Equal(err.Error(), mgo.ErrNotFound.Error())
}

func (suite *NetworkTestSuite) TestDeleteEmptyNetwork() {
	//Remove with non-exist network id
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/networks/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusNotFound, httpWriter)
}

//Fot Get/List, we only return mongo document
func (suite *NetworkTestSuite) TestGetNetwork() {
	tName := namesgenerator.GetRandomName(0)
	tType := "ovs"
	eth1 := entity.PhysicalPort{
		Name:     suite.ifName,
		MTU:      1500,
		VlanTags: []int{2043, 2143, 2243},
	}
	network := entity.Network{
		ID:            bson.NewObjectId(),
		BridgeName:    tName,
		BridgeType:    tType,
		NodeName:      suite.nodeName,
		PhysicalPorts: []entity.PhysicalPort{eth1},
	}

	//Create data into mongo manually
	suite.session.C(entity.NetworkCollectionName).Insert(network)
	defer suite.session.Remove(entity.NetworkCollectionName, "bridgeName", tName)

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/networks/"+network.ID.Hex(), nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	network = entity.Network{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &network)
	suite.NoError(err)
	suite.Equal(tName, network.BridgeName)
	suite.Equal(eth1, network.PhysicalPorts[0])
	suite.Equal(tType, network.BridgeType)
	suite.Equal(suite.nodeName, network.NodeName)
}

func (suite *NetworkTestSuite) TestGetNetworkWithInvalidID() {

	//Get data with non-exits ID
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/networks/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusNotFound, httpWriter)
}

func (suite *NetworkTestSuite) TestListNetwork() {
	networks := []entity.Network{}
	for i := 0; i < 3; i++ {
		networks = append(networks, entity.Network{
			BridgeName:    namesgenerator.GetRandomName(0),
			BridgeType:    "ovs",
			NodeName:      suite.nodeName,
			PhysicalPorts: []entity.PhysicalPort{}})
	}

	for _, v := range networks {
		suite.session.C(entity.NetworkCollectionName).Insert(v)
		defer suite.session.Remove(entity.NetworkCollectionName, "bridgeName", v.BridgeName)
	}

	//list data by default page and page_size
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/networks/", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	retNetworks := []entity.Network{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &retNetworks)
	suite.NoError(err)
	suite.Equal(len(networks), len(retNetworks))
	for i, v := range retNetworks {
		suite.Equal(networks[i].BridgeName, v.BridgeName)
		suite.Equal(networks[i].BridgeType, v.BridgeType)
		suite.Equal(networks[i].NodeName, v.NodeName)
		suite.Equal(networks[i].PhysicalPorts, v.PhysicalPorts)
	}

	//list data by different page and page_size
	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/networks?page=1&page_size=3", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	retNetworks = []entity.Network{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &retNetworks)
	suite.NoError(err)
	suite.Equal(len(networks), len(retNetworks))
	for i, v := range retNetworks {
		suite.Equal(networks[i].BridgeName, v.BridgeName)
		suite.Equal(networks[i].BridgeType, v.BridgeType)
		suite.Equal(networks[i].NodeName, v.NodeName)
		suite.Equal(networks[i].PhysicalPorts, v.PhysicalPorts)
	}

	//list data by different page and page_size
	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/networks?page=1&page_size=1", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	retNetworks = []entity.Network{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &retNetworks)
	suite.NoError(err)
	suite.Equal(1, len(retNetworks))
	for i, v := range retNetworks {
		suite.Equal(networks[i].BridgeName, v.BridgeName)
		suite.Equal(networks[i].BridgeType, v.BridgeType)
		suite.Equal(networks[i].NodeName, v.NodeName)
		suite.Equal(networks[i].PhysicalPorts, v.PhysicalPorts)
	}
}

func (suite *NetworkTestSuite) TestListNetworkWithInvalidPage() {
	//Get data with non-exits ID
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/networks?page=asdd", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/networks?page_size=asdd", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/networks?page=-1", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusInternalServerError, httpWriter)
}
