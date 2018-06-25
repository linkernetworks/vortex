package server

import (
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	fakeclientset "k8s.io/client-go/kubernetes/fake"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type NetworkTestSuite struct {
	suite.Suite
	kubectl    *kc.KubeCtl
	fakeclient *fakeclientset.Clientset
	wc         *restful.Container
	session    *mongo.Session
	ifName     string
	nodeName   string
}

func (suite *NetworkTestSuite) SetupTest() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	//init session
	suite.session = sp.Mongo.NewSession()
	//init restful container
	suite.wc = restful.NewContainer()
	service := newNetworkService(sp)
	suite.wc.Add(service)

	//init fakeclient
	suite.fakeclient = fakeclientset.NewSimpleClientset()
	namespace := "default"
	suite.kubectl = kc.New(suite.fakeclient, namespace)

	sp.KubeCtl = suite.kubectl
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
	_, err := suite.fakeclient.CoreV1().Nodes().Create(&node)
	assert.NoError(suite.T(), err)

	//There's a length limit of link name
	suite.ifName = namesgenerator.GetRandomName(0)[0:8]
	pName := namesgenerator.GetRandomName(0)[0:8]
	//Create a veth for testing
	err = exec.Command("ip", "link", "add", suite.ifName, "type", "veth", "peer", "name", pName).Run()
	assert.NoError(suite.T(), err)
}

func (suite *NetworkTestSuite) TearDownTest() {
	err := exec.Command("ip", "link", "del", suite.ifName).Run()
	assert.NoError(suite.T(), err)
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
	assert.NoError(suite.T(), err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/networks", bodyReader)
	assert.NoError(suite.T(), err)

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
	assert.NoError(suite.T(), err)
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
	assert.NoError(suite.T(), err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/networks", bodyReader)
	assert.NoError(suite.T(), err)

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
	defer suite.session.Remove(entity.NetworkCollectionName, "bridgeName", tName)

	httpRequestDelete, err := http.NewRequest("DELETE", "http://localhost:7890/v1/networks/"+network.ID.Hex(), nil)
	httpWriterDelete := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriterDelete, httpRequestDelete)
	assertResponseCode(suite.T(), http.StatusOK, httpWriterDelete)
	err = suite.session.FindOne(entity.NetworkCollectionName, bson.M{"_id": network.ID}, &network)
	assert.Equal(suite.T(), err.Error(), mgo.ErrNotFound.Error())
}

func (suite *NetworkTestSuite) TestDeleteEmptyNetwork() {
	//Remove with non-exist network id
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/networks/"+bson.NewObjectId().Hex(), nil)
	assert.NoError(suite.T(), err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusNotFound, httpWriter)
}

func TestGetNetwork(t *testing.T) {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	tName := namesgenerator.GetRandomName(0)
	tType := "ovs"
	tNodeName := namesgenerator.GetRandomName(0)
	eth1 := entity.PhysicalPort{
		Name:     "eth1",
		MTU:      1500,
		VlanTags: []int{2043, 2143, 2243},
	}
	network := entity.Network{
		ID:            bson.NewObjectId(),
		BridgeName:    tName,
		BridgeType:    tType,
		NodeName:      tNodeName,
		PhysicalPorts: []entity.PhysicalPort{eth1},
	}

	//Create data into mongo manually
	session := sp.Mongo.NewSession()
	defer session.Close()
	session.C(entity.NetworkCollectionName).Insert(network)
	defer session.Remove(entity.NetworkCollectionName, "bridgeName", tName)

	httpRequestGet, err := http.NewRequest("GET", "http://localhost:7890/v1/networks/"+network.ID.Hex(), nil)
	assert.NoError(t, err)

	httpWriterGet := httptest.NewRecorder()
	wcGet := restful.NewContainer()
	serviceGet := newNetworkService(sp)
	wcGet.Add(serviceGet)
	wcGet.Dispatch(httpWriterGet, httpRequestGet)
	assertResponseCode(t, http.StatusOK, httpWriterGet)

	network = entity.Network{}
	err = json.Unmarshal(httpWriterGet.Body.Bytes(), &network)
	assert.NoError(t, err)
	assert.Equal(t, tName, network.BridgeName)
	assert.Equal(t, eth1, network.PhysicalPorts[0])
	assert.Equal(t, tType, network.BridgeType)
	assert.Equal(t, tNodeName, network.NodeName)
}

func TestGetNetworkWithInvalidID(t *testing.T) {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	//Get data with non-exits ID
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/networks/"+bson.NewObjectId().Hex(), nil)
	assert.NoError(t, err)

	httpWriter := httptest.NewRecorder()
	wc := restful.NewContainer()
	service := newNetworkService(sp)
	wc.Add(service)
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusNotFound, httpWriter)
}

func TestListNetwork(t *testing.T) {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	networks := []entity.Network{}

	for i := 0; i < 3; i++ {
		networks = append(networks, entity.Network{
			BridgeName: namesgenerator.GetRandomName(0),
			BridgeType: "ovs",
			NodeName:   namesgenerator.GetRandomName(0),
			PhysicalPorts: []entity.PhysicalPort{
				{namesgenerator.GetRandomName(0), 1500, []int{1234, 123, 432}},
			}})
	}

	session := sp.Mongo.NewSession()
	defer session.Close()
	for _, v := range networks {
		session.C(entity.NetworkCollectionName).Insert(v)
		defer session.Remove(entity.NetworkCollectionName, "bridgeName", v.BridgeName)
	}

	//list data by default page and page_size
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/networks/", nil)
	assert.NoError(t, err)

	httpWriter := httptest.NewRecorder()
	wc := restful.NewContainer()
	service := newNetworkService(sp)
	wc.Add(service)
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusOK, httpWriter)

	retNetworks := []entity.Network{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &retNetworks)
	assert.NoError(t, err)
	assert.Equal(t, len(networks), len(retNetworks))
	for i, v := range retNetworks {
		assert.Equal(t, networks[i].BridgeName, v.BridgeName)
		assert.Equal(t, networks[i].BridgeType, v.BridgeType)
		assert.Equal(t, networks[i].NodeName, v.NodeName)
		assert.Equal(t, networks[i].PhysicalPorts, v.PhysicalPorts)
	}

	//list data by different page and page_size
	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/networks?page=1&page_size=3", nil)
	assert.NoError(t, err)

	httpWriter = httptest.NewRecorder()
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusOK, httpWriter)

	retNetworks = []entity.Network{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &retNetworks)
	assert.NoError(t, err)
	assert.Equal(t, len(networks), len(retNetworks))
	for i, v := range retNetworks {
		assert.Equal(t, networks[i].BridgeName, v.BridgeName)
		assert.Equal(t, networks[i].BridgeType, v.BridgeType)
		assert.Equal(t, networks[i].NodeName, v.NodeName)
		assert.Equal(t, networks[i].PhysicalPorts, v.PhysicalPorts)
	}

	//list data by different page and page_size
	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/networks?page=1&page_size=1", nil)
	assert.NoError(t, err)

	httpWriter = httptest.NewRecorder()
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusOK, httpWriter)

	retNetworks = []entity.Network{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &retNetworks)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(retNetworks))
	for i, v := range retNetworks {
		assert.Equal(t, networks[i].BridgeName, v.BridgeName)
		assert.Equal(t, networks[i].BridgeType, v.BridgeType)
		assert.Equal(t, networks[i].NodeName, v.NodeName)
		assert.Equal(t, networks[i].PhysicalPorts, v.PhysicalPorts)
	}
}

func TestListNetworkWithInvalidPage(t *testing.T) {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	//Get data with non-exits ID
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/networks?page=asdd", nil)
	assert.NoError(t, err)

	httpWriter := httptest.NewRecorder()
	wc := restful.NewContainer()
	service := newNetworkService(sp)
	wc.Add(service)
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/networks?page_size=asdd", nil)
	assert.NoError(t, err)

	httpWriter = httptest.NewRecorder()
	service = newNetworkService(sp)
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/networks?page=-1", nil)
	assert.NoError(t, err)

	httpWriter = httptest.NewRecorder()
	service = newNetworkService(sp)
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusInternalServerError, httpWriter)

}
