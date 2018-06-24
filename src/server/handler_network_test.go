package server

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/stretchr/testify/assert"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestCreateNetwork(t *testing.T) {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	eth1 := entity.PhysicalPort{
		Name:     "eth1",
		MTU:      1500,
		VlanTags: []int{2043, 2143, 2243},
	}

	tName := namesgenerator.GetRandomName(0)
	network := entity.Network{
		BridgeName:    tName,
		BridgeType:    "ovs",
		NodeName:      "create-network-node",
		PhysicalPorts: []entity.PhysicalPort{eth1},
	}

	session := sp.Mongo.NewSession()
	defer session.Close()

	bodyBytes, err := json.MarshalIndent(network, "", "  ")
	assert.NoError(t, err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/networks", bodyReader)
	assert.NoError(t, err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	wc := restful.NewContainer()
	service := newNetworkService(sp)
	wc.Add(service)
	wc.Dispatch(httpWriter, httpRequest)
	defer session.Remove(entity.NetworkCollectionName, "bridgeName", tName)

	//We use the new write but empty input
	httpWriter = httptest.NewRecorder()
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusBadRequest, httpWriter)
	//Create again and it should fail since the name exist
	bodyReader = strings.NewReader(string(bodyBytes))
	httpRequest, err = http.NewRequest("POST", "http://localhost:7890/v1/networks", bodyReader)
	assert.NoError(t, err)
	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter = httptest.NewRecorder()
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusConflict, httpWriter)
}

func TestWrongVlangTag(t *testing.T) {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	tName := namesgenerator.GetRandomName(0)
	network := entity.Network{
		BridgeName: tName,
		BridgeType: "ovs",
		NodeName:   "wron-vlan-node3",
		PhysicalPorts: []entity.PhysicalPort{
			{
				Name:     "eth1",
				MTU:      1500,
				VlanTags: []int{1234, 2143, 2243},
			},
			{
				Name:     "eth1",
				MTU:      1500,
				VlanTags: []int{1234, 2143, 50000},
			},
		}}

	session := sp.Mongo.NewSession()
	defer session.Close()

	bodyBytes, err := json.MarshalIndent(network, "", "  ")
	assert.NoError(t, err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/networks", bodyReader)
	assert.NoError(t, err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	wc := restful.NewContainer()
	service := newNetworkService(sp)
	wc.Add(service)
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusBadRequest, httpWriter)
}

func TestDeleteNetwork(t *testing.T) {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	tName := namesgenerator.GetRandomName(0)
	network := entity.Network{
		ID:            bson.NewObjectId(),
		BridgeName:    tName,
		BridgeType:    "ovs",
		NodeName:      "delete-network-node",
		PhysicalPorts: []entity.PhysicalPort{},
	}

	//Create data into mongo manually
	session := sp.Mongo.NewSession()
	defer session.Close()
	session.C(entity.NetworkCollectionName).Insert(network)
	defer session.Remove(entity.NetworkCollectionName, "bridgeName", tName)

	httpRequestDelete, err := http.NewRequest("DELETE", "http://localhost:7890/v1/networks/"+network.ID.Hex(), nil)
	httpWriterDelete := httptest.NewRecorder()
	wcDelete := restful.NewContainer()
	serviceDelete := newNetworkService(sp)
	wcDelete.Add(serviceDelete)
	wcDelete.Dispatch(httpWriterDelete, httpRequestDelete)
	assertResponseCode(t, http.StatusOK, httpWriterDelete)

	err = session.FindOne(entity.NetworkCollectionName, bson.M{"_id": network.ID}, &network)
	assert.Equal(t, err.Error(), mgo.ErrNotFound.Error())
}

func TestDeleteEmptyNetwork(t *testing.T) {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	//Remove with non-exist network id
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/networks/"+bson.NewObjectId().Hex(), nil)
	assert.NoError(t, err)

	httpWriter := httptest.NewRecorder()
	wc := restful.NewContainer()
	service := newNetworkService(sp)
	wc.Add(service)
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusNotFound, httpWriter)
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
}
