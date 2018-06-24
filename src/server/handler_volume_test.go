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
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2/bson"
	corev1 "k8s.io/api/core/v1"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type VolumeTestSuite struct {
	suite.Suite
	wc              *restful.Container
	session         *mongo.Session
	storageProvider entity.StorageProvider
}

func (suite *VolumeTestSuite) SetupTest() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	//init session
	suite.session = sp.Mongo.NewSession()
	//init restful container
	suite.wc = restful.NewContainer()
	service := newVolumeService(sp)
	suite.wc.Add(service)
	//init a StorageProvider
	suite.storageProvider = entity.StorageProvider{
		ID:          bson.NewObjectId(),
		Type:        "nfs",
		DisplayName: namesgenerator.GetRandomName(0),
	}
	err := suite.session.Insert(entity.StorageProviderCollectionName, suite.storageProvider)
	assert.NoError(suite.T(), err)
}

func (suite *VolumeTestSuite) TearDownTest() {
	suite.session.Remove(entity.StorageProviderCollectionName, "_id", suite.storageProvider.ID)
}

func TestVolumeSuite(t *testing.T) {
	suite.Run(t, new(VolumeTestSuite))
}

func (suite *VolumeTestSuite) TestCreateVolume() {
	tName := namesgenerator.GetRandomName(0)
	tAccessMode := corev1.PersistentVolumeAccessMode("ReadOnlyMany")
	tCapacity := "500G"
	volume := entity.Volume{
		Name:                tName,
		StorageProviderName: suite.storageProvider.DisplayName,
		Capacity:            tCapacity,
		AccessMode:          tAccessMode,
	}

	bodyBytes, err := json.MarshalIndent(volume, "", "  ")
	assert.NoError(suite.T(), err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/volume", bodyReader)
	assert.NoError(suite.T(), err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
	defer suite.session.Remove(entity.VolumeCollectionName, "name", volume.Name)

	//load data to check
	retVolume := entity.Volume{}
	err = suite.session.FindOne(entity.VolumeCollectionName, bson.M{"name": volume.Name}, &retVolume)
	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), "", retVolume.ID)
	assert.Equal(suite.T(), volume.Name, retVolume.Name)
	assert.Equal(suite.T(), volume.StorageProviderName, retVolume.StorageProviderName)
	assert.Equal(suite.T(), volume.AccessMode, retVolume.AccessMode)
	assert.Equal(suite.T(), volume.Capacity, retVolume.Capacity)
	assert.NotEqual(suite.T(), "", retVolume.MetaName)

	//We use the new write but empty input which will cause the readEntity Error
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
	//Create again and it should fail since the name exist
	bodyReader = strings.NewReader(string(bodyBytes))
	httpRequest, err = http.NewRequest("POST", "http://localhost:7890/v1/volume", bodyReader)
	assert.NoError(suite.T(), err)
	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusConflict, httpWriter)
}

func (suite *VolumeTestSuite) TestCreateVolumeWithInvalidParameter() {
	tName := namesgenerator.GetRandomName(0)
	tAccessMode := corev1.PersistentVolumeAccessMode("ReadOnlyMany")
	tCapacity := "500G"
	volume := entity.Volume{
		Name:                tName,
		StorageProviderName: namesgenerator.GetRandomName(0),
		Capacity:            tCapacity,
		AccessMode:          tAccessMode,
	}

	bodyBytes, err := json.MarshalIndent(volume, "", "  ")
	assert.NoError(suite.T(), err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/volume", bodyReader)
	assert.NoError(suite.T(), err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
	defer suite.session.Remove(entity.VolumeCollectionName, "name", volume.Name)

}
