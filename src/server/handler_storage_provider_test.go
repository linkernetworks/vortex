package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/stretchr/testify/assert"
)

func TestCreateStorageProvider(t *testing.T) {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	//Testing parameter
	tName := "hello world"
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
	defer session.Remove(entity.StorageProviderCollectionName, "displayName", tName)
	assertResponseCode(t, http.StatusOK, httpWriter)
}
