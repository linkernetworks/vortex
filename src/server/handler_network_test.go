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
		NodeName:      "node1",
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
	defer session.Remove(entity.NetworkCollectionName, "bridegName", tName)

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
		NodeName:   "node1",
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
		BridgeName:    tName,
		BridgeType:    "ovs",
		NodeName:      "node1",
		PhysicalPorts: []entity.PhysicalPort{},
	}

	//Create data into mongo manually
	session := sp.Mongo.NewSession()
	defer session.Close()
	session.C(entity.NetworkCollectionName).Insert(network)
	defer session.Remove(entity.NetworkCollectionName, "bridgeName", tName)

	//Reload the data to get the objectID
	network = entity.Network{}
	q := bson.M{"bridgeName": tName}
	err := session.FindOne(entity.NetworkCollectionName, q, &network)
	assert.NoError(t, err)

	httpRequestDelete, err := http.NewRequest("DELETE", "http://localhost:7890/v1/networks/"+network.ID.Hex(), nil)
	httpWriterDelete := httptest.NewRecorder()
	wcDelete := restful.NewContainer()
	serviceDelete := newNetworkService(sp)
	wcDelete.Add(serviceDelete)
	wcDelete.Dispatch(httpWriterDelete, httpRequestDelete)
	assertResponseCode(t, http.StatusOK, httpWriterDelete)

	err = session.FindOne(entity.NetworkCollectionName, q, &network)
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

	//Reload the data to get the objectID
	network = entity.Network{}
	q := bson.M{"bridgeName": tName}
	err := session.FindOne(entity.NetworkCollectionName, q, &network)
	assert.NoError(t, err)

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
