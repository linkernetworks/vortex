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

func (suite *StorageTestSuite) TearDownSuite() {}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}

func (suite *StorageTestSuite) TestCreateStorage() {
	//Testing parameter
	tName := namesgenerator.GetRandomName(0)
	storage := entity.Storage{
		Type:             entity.FakeStorageType,
		Name:             tName,
		StorageClassName: tName,
		IP:               "192.168.5.100",
		PATH:             "/myspace",
		Fake: &entity.FakeStorage{
			FakeParameter: "fake~",
		},
	}

	bodyBytes, err := json.MarshalIndent(storage, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/storage", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	defer suite.session.Remove(entity.StorageCollectionName, "name", tName)
	assertResponseCode(suite.T(), http.StatusCreated, httpWriter)
	//Empty data
	//We use the new write but empty input
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
	//Create again and it should fail since the name exist
	bodyReader = strings.NewReader(string(bodyBytes))
	httpRequest, err = http.NewRequest("POST", "http://localhost:7890/v1/storage", bodyReader)
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
			Type:             entity.FakeStorageType,
			Name:             namesgenerator.GetRandomName(0),
			StorageClassName: namesgenerator.GetRandomName(1),
			IP:               "192.168.5.100",
			PATH:             "/myspace",
			Fake: &entity.FakeStorage{
				FakeParameter: "",
			},
		},
			http.StatusBadRequest},
		{"CreateFail", entity.Storage{
			Name:             namesgenerator.GetRandomName(0),
			StorageClassName: namesgenerator.GetRandomName(1),
			Type:             entity.FakeStorageType,
			IP:               "192.168.5.100",
			PATH:             "/myspace",
			Fake: &entity.FakeStorage{
				FakeParameter: "Yo",
				IWantFail:     true,
			},
		},
			http.StatusInternalServerError},
		{"inValidNFSIP", entity.Storage{
			Type:             entity.FakeStorageType,
			Name:             namesgenerator.GetRandomName(1),
			StorageClassName: namesgenerator.GetRandomName(0),
			IP:               "256.256.256.256",
			PATH:             "/myspace",
			Fake: &entity.FakeStorage{
				FakeParameter: "Yo",
			},
		},
			http.StatusBadRequest},
		{"lackStorageName", entity.Storage{
			Type:             entity.FakeStorageType,
			StorageClassName: namesgenerator.GetRandomName(0),
			IP:               "256.256.256.256",
			PATH:             "/myspace",
			Fake: &entity.FakeStorage{
				FakeParameter: "Yo",
			},
		},
			http.StatusBadRequest},
		{"StorageTypeError", entity.Storage{
			Name:             namesgenerator.GetRandomName(0),
			StorageClassName: namesgenerator.GetRandomName(1),
			Type:             "none-exist",
			IP:               "192.168.5.100",
			PATH:             "/myspace",
			Fake: &entity.FakeStorage{
				FakeParameter: "Yo",
				IWantFail:     true,
			},
		},
			http.StatusBadRequest},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.cases, func(t *testing.T) {
			bodyBytes, err := json.MarshalIndent(tc.storage, "", "  ")
			suite.NoError(err)

			bodyReader := strings.NewReader(string(bodyBytes))
			httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/storage", bodyReader)
			suite.NoError(err)

			httpRequest.Header.Add("Content-Type", "application/json")
			httpWriter := httptest.NewRecorder()
			suite.wc.Dispatch(httpWriter, httpRequest)
			assertResponseCode(suite.T(), tc.errorCode, httpWriter)
		})
	}

}

func (suite *StorageTestSuite) TestDeleteStorage() {
	//Testing parameter
	tName := namesgenerator.GetRandomName(0)
	storage := entity.Storage{
		ID:   bson.NewObjectId(),
		Type: entity.FakeStorageType,
		Name: tName,
		Fake: &entity.FakeStorage{
			FakeParameter: "fake~",
		},
	}

	suite.session.C(entity.StorageCollectionName).Insert(storage)
	defer suite.session.Remove(entity.StorageCollectionName, "name", tName)

	bodyBytes, err := json.MarshalIndent(storage, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/storage/"+storage.ID.Hex(), bodyReader)
	suite.NoError(err)
	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *NetworkTestSuite) TestDeleteEmptyStorage() {
	//Remove with non-exist network id
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/storage/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusNotFound, httpWriter)
}

func (suite *StorageTestSuite) TestInValidDeleteStorage() {
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/storage/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusNotFound, httpWriter)
}

func (suite *StorageTestSuite) TestDeleteStorageFail() {
	testCases := []struct {
		cases     string
		storage   entity.Storage
		errorCode int
	}{
		{"DeleteStorage", entity.Storage{
			ID:   bson.NewObjectId(),
			Name: namesgenerator.GetRandomName(0),
			Type: entity.FakeStorageType,
			Fake: &entity.FakeStorage{
				FakeParameter: "Yo-Delete-Fail",
				IWantFail:     true,
			}},
			http.StatusInternalServerError},
		{"StorageTypeError", entity.Storage{
			ID:   bson.NewObjectId(),
			Name: namesgenerator.GetRandomName(0),
			Type: "non-exist",
			Fake: &entity.FakeStorage{
				FakeParameter: "Yo-Delete-Fail",
				IWantFail:     true,
			}},
			http.StatusBadRequest},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.cases, func(t *testing.T) {
			suite.session.C(entity.StorageCollectionName).Insert(tc.storage)
			defer suite.session.Remove(entity.StorageCollectionName, "name", tc.storage.Name)

			httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/storage/"+tc.storage.ID.Hex(), nil)
			suite.NoError(err)

			httpRequest.Header.Add("Content-Type", "application/json")
			httpWriter := httptest.NewRecorder()
			suite.wc.Dispatch(httpWriter, httpRequest)
			assertResponseCode(suite.T(), tc.errorCode, httpWriter)
		})
	}
}

func (suite *StorageTestSuite) TestListStorage() {
	storages := []entity.Storage{}

	count := 3
	for i := 0; i < count; i++ {
		storages = append(storages, entity.Storage{
			Name: namesgenerator.GetRandomName(0),
			Type: entity.FakeStorageType,
			Fake: &entity.FakeStorage{
				FakeParameter: "Yo",
				IWantFail:     false,
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

	for _, v := range storages {
		err := suite.session.C(entity.StorageCollectionName).Insert(v)
		defer suite.session.Remove(entity.StorageCollectionName, "name", v.Name)
		suite.NoError(err)
	}

	for _, tc := range testCases {
		caseName := "page:pageSize" + tc.page + ":" + tc.pageSize
		suite.T().Run(caseName, func(t *testing.T) {
			url := "http://localhost:7890/v1/storage/"
			if tc.page != "" || tc.pageSize != "" {
				url = "http://localhost:7890/v1/storage?"
				url += "page=" + tc.page + "%" + "page_size" + tc.pageSize
			}
			httpRequest, err := http.NewRequest("GET", url, nil)

			suite.NoError(err)

			httpWriter := httptest.NewRecorder()
			suite.wc.Dispatch(httpWriter, httpRequest)
			assertResponseCode(suite.T(), http.StatusOK, httpWriter)

			retStorages := []entity.Storage{}
			err = json.Unmarshal(httpWriter.Body.Bytes(), &retStorages)
			suite.NoError(err)
			suite.Equal(tc.expectSize, len(retStorages))
			for i, v := range retStorages {
				suite.Equal(storages[i].Name, v.Name)
				suite.Equal(storages[i].Type, v.Type)
			}
		})
	}
}

func (suite *StorageTestSuite) TestListInvalidStorage() {
	//Invliad page size
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/storage?page=0", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusInternalServerError, httpWriter)

	//Invliad page type
	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/storage?page=asd", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	//Invliad page_size type
	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/storage?page_size=asd", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
}
