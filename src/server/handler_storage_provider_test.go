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
	"gopkg.in/mgo.v2/bson"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type StorageProviderTestSuite struct {
	suite.Suite
	wc              *restful.Container
	storageProvider entity.NFSStorageProvider
	session         *mongo.Session
}

func (suite *StorageProviderTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	//init restful container
	suite.wc = restful.NewContainer()
	service := newStorageProviderService(sp)
	suite.wc.Add(service)

	//init session
	suite.session = sp.Mongo.NewSession()
}

func (suite *StorageProviderTestSuite) TearDownSuite() {
}

func TestStorageProviderSuite(t *testing.T) {
	suite.Run(t, new(StorageProviderTestSuite))
}

func (suite *StorageProviderTestSuite) TestCreateStorageProvider() {
	//Testing parameter
	tName := namesgenerator.GetRandomName(0)
	tType := "nfs"
	tIP := "1.2.3.4"
	tPath := "/exports"
	storageProvider := entity.StorageProvider{
		Type:        tType,
		DisplayName: tName,
		NFSStorageProvider: entity.NFSStorageProvider{
			IP:   tIP,
			PATH: tPath,
		},
	}

	bodyBytes, err := json.MarshalIndent(storageProvider, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/storageprovider", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	defer suite.session.Remove(entity.StorageProviderCollectionName, "displayName", tName)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
	//Empty data
	//We use the new write but empty input
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
	//Create again and it should fail since the name exist
	bodyReader = strings.NewReader(string(bodyBytes))
	httpRequest, err = http.NewRequest("POST", "http://localhost:7890/v1/storageprovider", bodyReader)
	suite.NoError(err)
	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusConflict, httpWriter)
}

func (suite *StorageProviderTestSuite) TestDeleteStorageProvider() {
	//Testing parameter
	tName := namesgenerator.GetRandomName(0)
	tType := "nfs"
	tIP := "1.2.3.4"
	tPath := "/exports"
	storageProvider := entity.StorageProvider{
		ID:          bson.NewObjectId(),
		Type:        tType,
		DisplayName: tName,
		NFSStorageProvider: entity.NFSStorageProvider{
			IP:   tIP,
			PATH: tPath,
		},
	}

	suite.session.C(entity.StorageProviderCollectionName).Insert(storageProvider)
	defer suite.session.Remove(entity.StorageProviderCollectionName, "displayName", tName)

	bodyBytes, err := json.MarshalIndent(suite.storageProvider, "", "  ")
	suite.NoError(err)

	//Create again and it should fail since the name exist
	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/storageprovider/"+storageProvider.ID.Hex(), bodyReader)
	suite.NoError(err)
	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *StorageProviderTestSuite) TestInValidDeleteStorageProvider() {
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/storageprovider/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusNotFound, httpWriter)
}

func (suite *StorageProviderTestSuite) TestListStorageProvider() {
	storageProviders := []entity.StorageProvider{}
	for i := 0; i < 3; i++ {
		storageProviders = append(storageProviders, entity.StorageProvider{
			ID:          bson.NewObjectId(),
			DisplayName: namesgenerator.GetRandomName(0),
			Type:        "nfs",
			NFSStorageProvider: entity.NFSStorageProvider{
				IP:   "1.2.3.4",
				PATH: "/expots",
			},
		})
	}

	for _, v := range storageProviders {
		err := suite.session.C(entity.StorageProviderCollectionName).Insert(v)
		suite.NoError(err)
		defer suite.session.Remove(entity.StorageProviderCollectionName, "_id", v.ID)
	}

	//default page & page_size
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/storageprovider/", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	retStorageProviders := []entity.StorageProvider{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &retStorageProviders)
	suite.NoError(err)
	suite.Equal(len(storageProviders), len(retStorageProviders))
	for i, v := range retStorageProviders {
		suite.Equal(storageProviders[i].ID, v.ID)
		suite.Equal(storageProviders[i].DisplayName, v.DisplayName)
		suite.Equal(storageProviders[i].Type, v.Type)
		suite.Equal(storageProviders[i].IP, v.IP)
		suite.Equal(storageProviders[i].PATH, v.PATH)
	}

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/storageprovider?page=1&page_size=30", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	retStorageProviders = []entity.StorageProvider{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &retStorageProviders)
	suite.NoError(err)
	suite.Equal(len(storageProviders), len(retStorageProviders))
	for i, v := range retStorageProviders {
		suite.Equal(storageProviders[i].ID, v.ID)
		suite.Equal(storageProviders[i].DisplayName, v.DisplayName)
		suite.Equal(storageProviders[i].Type, v.Type)
		suite.Equal(storageProviders[i].IP, v.IP)
		suite.Equal(storageProviders[i].PATH, v.PATH)
	}
}

func (suite *StorageProviderTestSuite) TestListInvalidStorageProvider() {
	//Invliad page size
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/storageprovider?page=0", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusInternalServerError, httpWriter)

	//Invliad page type
	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/storageprovider?page=asd", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	//Invliad page_size type
	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/storageprovider?page_size=asd", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
}
