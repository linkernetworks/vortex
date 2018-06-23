package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/docker/docker/pkg/namesgenerator"
	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/stretchr/testify/assert"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

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
		Name:          tName,
		BridgeType:    "ovs",
		Node:          "node1",
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
	defer session.Remove(entity.NetworkCollectionName, "name", tName)

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
		Name:       tName,
		BridgeType: "ovs",
		Node:       "node1",
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
		Name:          tName,
		BridgeType:    "ovs",
		Node:          "node1",
		PhysicalPorts: []entity.PhysicalPort{},
	}

	//Create data into mongo manually
	session := sp.Mongo.NewSession()
	defer session.Close()
	session.C(entity.NetworkCollectionName).Insert(network)
	defer session.Remove(entity.NetworkCollectionName, "name", tName)

	//Reload the data to get the objectID

	network = entity.Network{}
	q := bson.M{"name": tName}
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

	//Remove non-exist network id
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/networks/"+bson.NewObjectId().Hex(), nil)
	assert.NoError(t, err)

	httpWriter := httptest.NewRecorder()
	wc := restful.NewContainer()
	service := newNetworkService(sp)
	wc.Add(service)
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusNotFound, httpWriter)
}

func TestUpdateNetwork(t *testing.T) {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	t.Skip()
	eth1 := entity.PhysicalPort{
		Name:     "eth1",
		MTU:      1500,
		VlanTags: []int{2043, 2143, 2243},
	}
	network := entity.Network{
		Name:          "ovsbr1",
		BridgeType:    "ovs",
		Node:          "node1",
		PhysicalPorts: []entity.PhysicalPort{eth1},
	}

	session := sp.Mongo.NewSession()
	defer session.RemoveAll(entity.NetworkCollectionName)

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
	assertResponseCode(t, http.StatusOK, httpWriter)

	updatedNetwork := entity.Network{
		Name: "Test",
	}

	bodyBytesUpdate, err := json.MarshalIndent(updatedNetwork, "", "  ")
	assert.NoError(t, err)

	network = entity.Network{}
	q := bson.M{"name": "ovsbr1"}
	err = session.FindOne(entity.NetworkCollectionName, q, &network)
	assert.NoError(t, err)

	bodyReaderUpdate := strings.NewReader(string(bodyBytesUpdate))
	httpRequestUpdate, err := http.NewRequest("PUT", "http://localhost:7890/v1/networks/"+network.ID.Hex(), bodyReaderUpdate)
	assert.NoError(t, err)

	httpRequestUpdate.Header.Add("Content-Type", "application/json")
	httpWriterUpdate := httptest.NewRecorder()
	wcUpdate := restful.NewContainer()
	serviceUpdate := newNetworkService(sp)
	wcUpdate.Add(serviceUpdate)
	wcUpdate.Dispatch(httpWriterUpdate, httpRequestUpdate)
	assertResponseCode(t, http.StatusOK, httpWriterUpdate)
}

func TestWrongUpdateNetwork(t *testing.T) {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	t.Skip()
	eth1 := entity.PhysicalPort{
		Name:     "eth1",
		MTU:      1500,
		VlanTags: []int{2043, 2143, 2243},
	}
	network := entity.Network{
		Name:          "ovsbr1",
		BridgeType:    "ovs",
		Node:          "node1",
		PhysicalPorts: []entity.PhysicalPort{eth1},
	}

	session := sp.Mongo.NewSession()
	defer session.RemoveAll(entity.NetworkCollectionName)

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
	assertResponseCode(t, http.StatusOK, httpWriter)

	updatedNetwork := entity.Network{
		Name: "obsbr2",
	}

	bodyBytesUpdate, err := json.MarshalIndent(updatedNetwork, "", "  ")
	assert.NoError(t, err)

	network = entity.Network{}
	q := bson.M{"name": "ovsbr1"}
	err = session.FindOne(entity.NetworkCollectionName, q, &network)
	assert.NoError(t, err)

	bodyReaderUpdate := strings.NewReader(string(bodyBytesUpdate))
	httpRequestUpdate, err := http.NewRequest("PUT", "http://localhost:7890/v1/networks/"+network.ID.Hex(), bodyReaderUpdate)
	assert.NoError(t, err)

	httpRequestUpdate.Header.Add("Content-Type", "application/json")
	httpWriterUpdate := httptest.NewRecorder()
	wcUpdate := restful.NewContainer()
	serviceUpdate := newNetworkService(sp)
	wcUpdate.Add(serviceUpdate)
	wcUpdate.Dispatch(httpWriterUpdate, httpRequestUpdate)
	assertResponseCode(t, http.StatusBadRequest, httpWriterUpdate)
}

func TestGetNetwork(t *testing.T) {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	t.Skip()
	eth1 := entity.PhysicalPort{
		Name:     "eth1",
		MTU:      1500,
		VlanTags: []int{2043, 2143, 2243},
	}
	network := entity.Network{
		Name:          "ovsbr1",
		BridgeType:    "ovs",
		Node:          "node1",
		PhysicalPorts: []entity.PhysicalPort{eth1},
	}
	session := sp.Mongo.NewSession()
	defer session.RemoveAll(entity.NetworkCollectionName)

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
	assertResponseCode(t, http.StatusOK, httpWriter)

	network = entity.Network{}
	q := bson.M{"name": "ovsbr1"}
	err = session.FindOne(entity.NetworkCollectionName, q, &network)
	assert.NoError(t, err)

	httpRequestGet, err := http.NewRequest("GET", "http://localhost:7890/v1/networks/"+network.ID.Hex(), nil)
	assert.NoError(t, err)

	httpWriterGet := httptest.NewRecorder()
	wcGet := restful.NewContainer()
	serviceGet := newNetworkService(sp)
	wcGet.Add(serviceGet)
	wcGet.Dispatch(httpWriterGet, httpRequestGet)
	assertResponseCode(t, http.StatusOK, httpWriterGet)
}
