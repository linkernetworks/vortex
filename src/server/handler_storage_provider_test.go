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
	//	"gopkg.in/mgo.v2/bson"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type StorageSuite struct {
	suite.Suite
	wc      *restful.Container
	session *mongo.Session
}

type StorageTestSuite struct {
	suite.Suite
	wc      *restful.Container
	session *mongo.Session
}

func (suite *StorageTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	//init restful container
	suite.wc = restful.NewContainer()
	service := newStorageService(sp)
	suite.wc.Add(service)

	//init session
	suite.session = sp.Mongo.NewSession()
}

func (suite *StorageTestSuite) TearDownSuite() {
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}

func (suite *StorageTestSuite) TestCreateStorage() {
	//Testing parameter
	tName := namesgenerator.GetRandomName(0)
	storage := entity.Storage{
		Type:        entity.FakeStorageType,
		DisplayName: tName,
		Fake: entity.FakeStorage{
			FakeParameter: "fake~",
		},
	}

	bodyBytes, err := json.MarshalIndent(storage, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/storageprovider", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	defer suite.session.Remove(entity.StorageCollectionName, "displayName", tName)
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

func (suite *StorageTestSuite) TestCreateStorageFail() {
	testCases := []struct {
		cases     string
		storage   entity.Storage
		errorCode int
	}{
		{"InvalidParameter", entity.Storage{
			DisplayName: namesgenerator.GetRandomName(0),
			Type:        entity.FakeStorageType,
			Fake: entity.FakeStorage{
				FakeParameter: "",
			}},
			http.StatusBadRequest},
		{"CreateFail", entity.Storage{
			DisplayName: namesgenerator.GetRandomName(0),
			Type:        entity.FakeStorageType,
			Fake: entity.FakeStorage{
				FakeParameter: "Yo",
				IWantFail:     true,
			}},
			http.StatusInternalServerError},
		{"StorageTypeError", entity.Storage{
			DisplayName: namesgenerator.GetRandomName(0),
			Type:        "non-exist",
			Fake: entity.FakeStorage{
				FakeParameter: "Yo",
				IWantFail:     true,
			}},
			http.StatusBadRequest},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.cases, func(t *testing.T) {
			bodyBytes, err := json.MarshalIndent(tc.storage, "", "  ")
			suite.NoError(err)

			bodyReader := strings.NewReader(string(bodyBytes))
			httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/storageprovider", bodyReader)
			suite.NoError(err)

			httpRequest.Header.Add("Content-Type", "application/json")
			httpWriter := httptest.NewRecorder()
			suite.wc.Dispatch(httpWriter, httpRequest)
			assertResponseCode(suite.T(), tc.errorCode, httpWriter)
		})
	}

}

/*
func (suite *StorageTestSuite) TestDeleteStorage() {
	//Testing parameter
	tName := namesgenerator.GetRandomName(0)
	tType := "nfs"
	tIP := "1.2.3.4"
	tPath := "/exports"
	storage := entity.Storage{
		ID:          bson.NewObjectId(),
		Type:        tType,
		DisplayName: tName,
		NFSStorageSetting: entity.NFSStorageSetting{
			IP:   tIP,
			PATH: tPath,
		},
	}

	suite.session.C(entity.StorageCollectionName).Insert(storage)
	defer suite.session.Remove(entity.StorageCollectionName, "displayName", tName)

	bodyBytes, err := json.MarshalIndent(suite.storage, "", "  ")
	suite.NoError(err)

	//Create again and it should fail since the name exist
	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/storageprovider/"+storage.ID.Hex(), bodyReader)
	suite.NoError(err)
	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *StorageTestSuite) TestInValidDeleteStorage() {
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/storageprovider/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusNotFound, httpWriter)
}

func (suite *StorageTestSuite) TestListStorage() {
	storages := []entity.Storage{}
	for i := 0; i < 3; i++ {
		storages = append(storages, entity.Storage{
			ID:          bson.NewObjectId(),
			DisplayName: namesgenerator.GetRandomName(0),
			Type:        "nfs",
			NFSStorageSetting: entity.NFSStorageSetting{
				IP:   "1.2.3.4",
				PATH: "/expots",
			},
		})
	}

	for _, v := range storages {
		err := suite.session.C(entity.StorageCollectionName).Insert(v)
		suite.NoError(err)
		defer suite.session.Remove(entity.StorageCollectionName, "_id", v.ID)
	}

	//default page & page_size
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/storageprovider/", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	retStorages := []entity.Storage{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &retStorages)
	suite.NoError(err)
	suite.Equal(len(storages), len(retStorages))
	for i, v := range retStorages {
		suite.Equal(storages[i].ID, v.ID)
		suite.Equal(storages[i].DisplayName, v.DisplayName)
		suite.Equal(storages[i].Type, v.Type)
		suite.Equal(storages[i].IP, v.IP)
		suite.Equal(storages[i].PATH, v.PATH)
	}

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/storageprovider?page=1&page_size=30", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	retStorages = []entity.Storage{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &retStorages)
	suite.NoError(err)
	suite.Equal(len(storages), len(retStorages))
	for i, v := range retStorages {
		suite.Equal(storages[i].ID, v.ID)
		suite.Equal(storages[i].DisplayName, v.DisplayName)
		suite.Equal(storages[i].Type, v.Type)
		suite.Equal(storages[i].IP, v.IP)
		suite.Equal(storages[i].PATH, v.PATH)
	}
}

func (suite *StorageTestSuite) TestListInvalidStorage() {
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
*/
