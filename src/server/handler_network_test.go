package server

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"

	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	corev1 "k8s.io/api/core/v1"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type NetworkTestSuite struct {
	suite.Suite
	wc      *restful.Container
	session *mongo.Session
	sp      *serviceprovider.Container
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

	suite.sp = sp
}

func (suite *NetworkTestSuite) TearDownSuite() {}

func TestNetworkSuite(t *testing.T) {
	suite.Run(t, new(NetworkTestSuite))
}

func (suite *NetworkTestSuite) TestCreateNetwork() {
	tName := namesgenerator.GetRandomName(0)
	network := entity.Network{
		Type:       entity.FakeNetworkType,
		IsDPDKPort: true, //for fake network, true means success,
		Name:       tName,
		VLANTags:   []int32{},
		BridgeName: namesgenerator.GetRandomName(0),
		Nodes: []entity.Node{
			entity.Node{
				Name:          tName,
				PhyInterfaces: []entity.PhyInterface{},
			},
		},
		CreatedAt: &time.Time{},
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
	defer suite.session.Remove(entity.NetworkCollectionName, "name", tName)

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

func (suite *NetworkTestSuite) TestCreateNetworkFail() {
	testCases := []struct {
		cases     string
		network   entity.Network
		errorCode int
	}{
		{"CreateFail",
			entity.Network{
				Type:       entity.FakeNetworkType,
				Name:       namesgenerator.GetRandomName(0),
				VLANTags:   []int32{},
				BridgeName: namesgenerator.GetRandomName(0),
				Nodes: []entity.Node{
					entity.Node{
						Name:          namesgenerator.GetRandomName(0),
						PhyInterfaces: []entity.PhyInterface{},
					},
				},
			},
			http.StatusInternalServerError},
		{"NetworkTypeError",
			entity.Network{
				Type:       "none-exist",
				Name:       namesgenerator.GetRandomName(0),
				VLANTags:   []int32{},
				BridgeName: namesgenerator.GetRandomName(0),
				Nodes: []entity.Node{
					entity.Node{
						Name:          namesgenerator.GetRandomName(0),
						PhyInterfaces: []entity.PhyInterface{},
					},
				},
			},
			http.StatusBadRequest},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.cases, func(t *testing.T) {
			bodyBytes, err := json.MarshalIndent(tc.network, "", "  ")
			suite.NoError(err)

			bodyReader := strings.NewReader(string(bodyBytes))
			httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/networks", bodyReader)
			suite.NoError(err)

			httpRequest.Header.Add("Content-Type", "application/json")
			httpWriter := httptest.NewRecorder()
			suite.wc.Dispatch(httpWriter, httpRequest)
			assertResponseCode(suite.T(), tc.errorCode, httpWriter)
		})
	}

}

func (suite *NetworkTestSuite) TestDeleteNetwork() {
	tName := namesgenerator.GetRandomName(0)
	network := entity.Network{
		ID:         bson.NewObjectId(),
		IsDPDKPort: true, //for fake network, true means success,
		Name:       tName,
		VLANTags:   []int32{},
		Type:       entity.FakeNetworkType,
		BridgeName: namesgenerator.GetRandomName(0),
		Nodes: []entity.Node{
			entity.Node{
				Name:          namesgenerator.GetRandomName(0),
				PhyInterfaces: []entity.PhyInterface{},
			},
		},
	}

	//Create data into mongo manually
	suite.session.C(entity.NetworkCollectionName).Insert(network)
	defer suite.session.Remove(entity.NetworkCollectionName, "name", tName)

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

func (suite *NetworkTestSuite) TestDeleteNetworkFail() {
	testCases := []struct {
		cases     string
		network   entity.Network
		errorCode int
	}{
		{"NetworkDeleteFail", entity.Network{
			ID:         bson.NewObjectId(),
			Type:       entity.FakeNetworkType,
			Name:       namesgenerator.GetRandomName(0),
			VLANTags:   []int32{},
			BridgeName: namesgenerator.GetRandomName(0),
			Nodes: []entity.Node{
				entity.Node{
					Name:          namesgenerator.GetRandomName(0),
					PhyInterfaces: []entity.PhyInterface{},
				},
			},
		},
			http.StatusInternalServerError},
		{"NetworkTypeError",
			entity.Network{
				ID:         bson.NewObjectId(),
				Name:       namesgenerator.GetRandomName(0),
				VLANTags:   []int32{},
				BridgeName: namesgenerator.GetRandomName(0),
				Type:       "none-exist",
				Nodes: []entity.Node{
					entity.Node{
						Name:          namesgenerator.GetRandomName(0),
						PhyInterfaces: []entity.PhyInterface{},
					},
				},
			},
			http.StatusBadRequest},
		{"PodStillUse", entity.Network{
			ID:         bson.NewObjectId(),
			Type:       entity.FakeNetworkType,
			Name:       namesgenerator.GetRandomName(0),
			VLANTags:   []int32{},
			BridgeName: namesgenerator.GetRandomName(0),
			Nodes: []entity.Node{
				entity.Node{
					Name:          namesgenerator.GetRandomName(0),
					PhyInterfaces: []entity.PhyInterface{},
				},
			},
		},
			http.StatusMethodNotAllowed},
	}

	//Create the Pod using the network.
	pod := entity.Pod{
		ID: bson.NewObjectId(),
		Networks: []entity.PodNetwork{
			{
				Name: testCases[2].network.Name,
			},
		},
	}
	suite.session.Insert(entity.PodCollectionName, pod)
	defer suite.session.Remove(entity.PodCollectionName, "_id", pod.ID)

	k8sPod := corev1.Pod{
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		}}
	suite.sp.KubeCtl.CreatePod(&k8sPod, "default")

	for _, tc := range testCases {
		suite.T().Run(tc.cases, func(t *testing.T) {
			suite.session.C(entity.NetworkCollectionName).Insert(tc.network)
			defer suite.session.Remove(entity.NetworkCollectionName, "name", tc.network.Name)

			httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/networks/"+tc.network.ID.Hex(), nil)
			suite.NoError(err)

			httpRequest.Header.Add("Content-Type", "application/json")
			httpWriter := httptest.NewRecorder()
			suite.wc.Dispatch(httpWriter, httpRequest)
			assertResponseCode(suite.T(), tc.errorCode, httpWriter)
		})
	}
}

//For Get/List, we only return mongo document
func (suite *NetworkTestSuite) TestGetNetwork() {
	tName := namesgenerator.GetRandomName(0)
	tType := entity.FakeNetworkType
	network := entity.Network{
		ID:       bson.NewObjectId(),
		Name:     tName,
		VLANTags: []int32{},
		Type:     tType,
		Nodes: []entity.Node{
			entity.Node{
				Name:          namesgenerator.GetRandomName(0),
				PhyInterfaces: []entity.PhyInterface{},
			},
		},
	}
	//Create data into mongo manually
	suite.session.C(entity.NetworkCollectionName).Insert(network)
	defer suite.session.Remove(entity.NetworkCollectionName, "name", tName)

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/networks/"+network.ID.Hex(), nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	network = entity.Network{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &network)
	suite.NoError(err)
	suite.Equal(tName, network.Name)
	suite.Equal(tType, network.Type)
}

func (suite *NetworkTestSuite) TestGetNetworkWithInvalidID() {
	// Get data with none-exits ID
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/networks/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusNotFound, httpWriter)
}

func (suite *NetworkTestSuite) TestGetNetworkStatus() {
	tName := namesgenerator.GetRandomName(0)
	tType := entity.FakeNetworkType
	network := entity.Network{
		ID:       bson.NewObjectId(),
		Name:     tName,
		VLANTags: []int32{},
		Type:     tType,
		Nodes: []entity.Node{
			entity.Node{
				Name:          namesgenerator.GetRandomName(0),
				PhyInterfaces: []entity.PhyInterface{},
			},
		},
	}
	//Create data into mongo manually
	suite.session.C(entity.NetworkCollectionName).Insert(network)
	defer suite.session.Remove(entity.NetworkCollectionName, "name", tName)

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/networks/status/"+network.ID.Hex(), nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	nameList := []string{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &nameList)
	suite.NoError(err)
}

func (suite *NetworkTestSuite) TestListNetwork() {
	networks := []entity.Network{}

	count := 3
	for i := 0; i < count; i++ {
		networks = append(networks, entity.Network{
			Type:       entity.FakeNetworkType,
			Name:       namesgenerator.GetRandomName(0),
			VLANTags:   []int32{},
			BridgeName: namesgenerator.GetRandomName(0),
			Nodes: []entity.Node{
				entity.Node{
					Name:          namesgenerator.GetRandomName(0),
					PhyInterfaces: []entity.PhyInterface{},
				},
			},
		})
	}

	testCases := []struct {
		page       string
		pageSize   string
		expectSize int
	}{
		{"", "", count},
		{"1", "1", count},
		{"1", "3", count},
	}

	for _, v := range networks {
		err := suite.session.C(entity.NetworkCollectionName).Insert(v)
		defer suite.session.Remove(entity.NetworkCollectionName, "name", v.Name)
		suite.NoError(err)
	}

	for _, tc := range testCases {
		caseName := "page:pageSize" + tc.page + ":" + tc.pageSize
		suite.T().Run(caseName, func(t *testing.T) {
			url := "http://localhost:7890/v1/networks/"
			if tc.page != "" || tc.pageSize != "" {
				url = "http://localhost:7890/v1/networks?"
				url += "page=" + tc.page + "%" + "page_size" + tc.pageSize
			}
			httpRequest, err := http.NewRequest("GET", url, nil)

			suite.NoError(err)

			httpWriter := httptest.NewRecorder()
			suite.wc.Dispatch(httpWriter, httpRequest)
			assertResponseCode(suite.T(), http.StatusOK, httpWriter)

			retNetworks := []entity.Network{}
			err = json.Unmarshal(httpWriter.Body.Bytes(), &retNetworks)
			suite.NoError(err)
			suite.Equal(tc.expectSize, len(retNetworks))
			for i, v := range retNetworks {
				suite.Equal(networks[i].Name, v.Name)
				suite.Equal(networks[i].Type, v.Type)
				suite.Equal(networks[i].Nodes[0].Name, v.Nodes[0].Name)
			}
		})
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
