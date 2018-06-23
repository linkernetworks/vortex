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
	"gopkg.in/mgo.v2/bson"
)

func TestCreateStorageProvider(t *testing.T) {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

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
	session := sp.Mongo.NewSession()

	bodyBytes, err := json.MarshalIndent(storageProvider, "", "  ")
	assert.NoError(t, err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/storageprovider", bodyReader)
	assert.NoError(t, err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	wc := restful.NewContainer()
	service := newStorageProviderService(sp)
	wc.Add(service)
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusOK, httpWriter)
	//Empty data
	//We use the new write but empty input
	httpWriter = httptest.NewRecorder()
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusBadRequest, httpWriter)
	//Create again and it should fail since the name exist
	bodyReader = strings.NewReader(string(bodyBytes))
	httpRequest, err = http.NewRequest("POST", "http://localhost:7890/v1/storageprovider", bodyReader)
	assert.NoError(t, err)
	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter = httptest.NewRecorder()
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusConflict, httpWriter)
	defer session.Remove(entity.StorageProviderCollectionName, "displayName", tName)
}

func TestListStorageProvider(t *testing.T) {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	session := sp.Mongo.NewSession()
	defer session.Close()
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
		err := session.C(entity.StorageProviderCollectionName).Insert(v)
		assert.NoError(t, err)
		defer session.Remove(entity.StorageProviderCollectionName, "_id", v.ID)
	}

	//default page & page_size
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/storageprovider/", nil)
	assert.NoError(t, err)

	httpWriter := httptest.NewRecorder()
	wc := restful.NewContainer()
	service := newStorageProviderService(sp)
	wc.Add(service)
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusOK, httpWriter)

	retStorageProviders := []entity.StorageProvider{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &retStorageProviders)
	assert.NoError(t, err)
	assert.Equal(t, len(storageProviders), len(retStorageProviders))
	for i, v := range retStorageProviders {
		assert.Equal(t, storageProviders[i].ID, v.ID)
		assert.Equal(t, storageProviders[i].DisplayName, v.DisplayName)
		assert.Equal(t, storageProviders[i].Type, v.Type)
		assert.Equal(t, storageProviders[i].IP, v.IP)
		assert.Equal(t, storageProviders[i].PATH, v.PATH)
	}

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/storageprovider?page=1&page_size=30", nil)
	assert.NoError(t, err)

	httpWriter = httptest.NewRecorder()
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusOK, httpWriter)

	retStorageProviders = []entity.StorageProvider{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &retStorageProviders)
	assert.NoError(t, err)
	assert.Equal(t, len(storageProviders), len(retStorageProviders))
	for i, v := range retStorageProviders {
		assert.Equal(t, storageProviders[i].ID, v.ID)
		assert.Equal(t, storageProviders[i].DisplayName, v.DisplayName)
		assert.Equal(t, storageProviders[i].Type, v.Type)
		assert.Equal(t, storageProviders[i].IP, v.IP)
		assert.Equal(t, storageProviders[i].PATH, v.PATH)
	}
}

func TestListInvalidStorageProvider(t *testing.T) {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	//Invliad page size
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/storageprovider?page=0", nil)
	assert.NoError(t, err)

	httpWriter := httptest.NewRecorder()
	wc := restful.NewContainer()
	service := newStorageProviderService(sp)
	wc.Add(service)
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusInternalServerError, httpWriter)

	//Invliad page type
	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/storageprovider?page=asd", nil)
	assert.NoError(t, err)

	httpWriter = httptest.NewRecorder()
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusBadRequest, httpWriter)

	//Invliad page_size type
	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/storageprovider?page_size=asd", nil)
	assert.NoError(t, err)

	httpWriter = httptest.NewRecorder()
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, http.StatusBadRequest, httpWriter)

}
