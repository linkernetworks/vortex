package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bitbucket.org/linkernetworks/vortex/src/entity"
	"bitbucket.org/linkernetworks/vortex/src/serviceprovider"
	restful "github.com/emicklei/go-restful"
	"github.com/influxdata/influxdb/pkg/testing/assert"
	"github.com/linkernetworks/config"
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
	network := entity.Network{
		DisplayName:   "OVS Bridge",
		BridgeName:    "obsbr1",
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
}

func TestWrongVlangTag(t *testing.T) {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	eth1 := entity.PhysicalPort{
		Name:     "eth1",
		MTU:      1500,
		VlanTags: []int{1234, 2143, 2243},
	}
	eth2 := entity.PhysicalPort{
		Name:     "eth2",
		MTU:      1500,
		VlanTags: []int{1234, 2143, 50000},
	}
	network := entity.Network{
		DisplayName:   "OVS Bridge",
		BridgeName:    "obsbr1",
		BridgeType:    "ovs",
		Node:          "node1",
		PhysicalPorts: []entity.PhysicalPort{eth1, eth2},
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
	assertResponseCode(t, http.StatusBadRequest, httpWriter)
}

func TestConflictDisplayName(t *testing.T) {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	eth1 := entity.PhysicalPort{
		Name:     "eth1",
		MTU:      1500,
		VlanTags: []int{2043, 2143, 2243},
	}
	network := entity.Network{
		DisplayName:   "OVS Bridge",
		BridgeName:    "obsbr1",
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

	networkWithSameInterface := entity.Network{
		DisplayName:   "OVS Bridge",
		BridgeName:    "obsbr2",
		BridgeType:    "ovs",
		Node:          "node1",
		PhysicalPorts: []entity.PhysicalPort{eth1},
	}

	bodyBytesConflict, err := json.MarshalIndent(networkWithSameInterface, "", "  ")
	assert.NoError(t, err)

	bodyReaderConflict := strings.NewReader(string(bodyBytesConflict))
	httpRequestConflict, err := http.NewRequest("POST", "http://localhost:7890/v1/networks", bodyReaderConflict)
	assert.NoError(t, err)

	httpRequestConflict.Header.Add("Content-Type", "application/json")
	httpWriterConflict := httptest.NewRecorder()
	wc.Dispatch(httpWriterConflict, httpRequestConflict)
	assertResponseCode(t, http.StatusConflict, httpWriterConflict)
}

func TestConflictBridgeName(t *testing.T) {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	eth1 := entity.PhysicalPort{
		Name:     "eth1",
		MTU:      1500,
		VlanTags: []int{2043, 2143, 2243},
	}
	network := entity.Network{
		DisplayName:   "OVS Bridge",
		BridgeName:    "obsbr1",
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

	networkWithSameName := entity.Network{
		DisplayName:   "OVS",
		BridgeName:    "obsbr1",
		BridgeType:    "ovs",
		Node:          "node2",
		PhysicalPorts: []entity.PhysicalPort{eth1},
	}

	bodyBytesConflict, err := json.MarshalIndent(networkWithSameName, "", "  ")
	assert.NoError(t, err)

	bodyReaderConflict := strings.NewReader(string(bodyBytesConflict))
	httpRequestConflict, err := http.NewRequest("POST", "http://localhost:7890/v1/networks", bodyReaderConflict)
	assert.NoError(t, err)

	httpRequestConflict.Header.Add("Content-Type", "application/json")
	httpWriterConflict := httptest.NewRecorder()
	wc.Dispatch(httpWriterConflict, httpRequestConflict)
	assertResponseCode(t, http.StatusConflict, httpWriterConflict)
}

func TestDeleteNetwork(t *testing.T) {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	eth1 := entity.PhysicalPort{
		Name:     "eth1",
		MTU:      1500,
		VlanTags: []int{2043, 2143, 2243},
	}
	network := entity.Network{
		DisplayName:   "OVS Bridge",
		BridgeName:    "obsbr1",
		BridgeType:    "ovs",
		Node:          "node1",
		PhysicalPorts: []entity.PhysicalPort{eth1},
	}

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

	session := sp.Mongo.NewSession()
	defer session.Close()

	network = entity.Network{}
	q := bson.M{"displayName": "OVS Bridge"}
	err = session.FindOne(entity.NetworkCollectionName, q, &network)
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

	eth1 := entity.PhysicalPort{
		Name:     "eth1",
		MTU:      1500,
		VlanTags: []int{2043, 2143, 2243},
	}
	network := entity.Network{
		DisplayName:   "OVS Bridge",
		BridgeName:    "obsbr1",
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
		DisplayName: "Test",
	}

	bodyBytesUpdate, err := json.MarshalIndent(updatedNetwork, "", "  ")
	assert.NoError(t, err)

	network = entity.Network{}
	q := bson.M{"displayName": "OVS Bridge"}
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

	eth1 := entity.PhysicalPort{
		Name:     "eth1",
		MTU:      1500,
		VlanTags: []int{2043, 2143, 2243},
	}
	network := entity.Network{
		DisplayName:   "OVS Bridge",
		BridgeName:    "obsbr1",
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
		DisplayName: "Test",
		BridgeName:  "obsbr2",
	}

	bodyBytesUpdate, err := json.MarshalIndent(updatedNetwork, "", "  ")
	assert.NoError(t, err)

	network = entity.Network{}
	q := bson.M{"displayName": "OVS Bridge"}
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
