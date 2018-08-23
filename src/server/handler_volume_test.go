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
	v "github.com/linkernetworks/vortex/src/volume"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2/bson"
	corev1 "k8s.io/api/core/v1"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type VolumeTestSuite struct {
	suite.Suite
	sp        *serviceprovider.Container
	wc        *restful.Container
	session   *mongo.Session
	storage   entity.Storage
	JWTBearer string
}

func (suite *VolumeTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	suite.sp = sp
	//init session
	suite.session = sp.Mongo.NewSession()
	//init restful container
	suite.wc = restful.NewContainer()
	service := newVolumeService(suite.sp)
	suite.wc.Add(service)

	token, _ := loginGetToken(suite.sp, suite.wc)
	suite.JWTBearer = "Bearer " + token

	//init a Storage
	suite.storage = entity.Storage{
		ID:   bson.NewObjectId(),
		Type: "nfs",
		Name: namesgenerator.GetRandomName(0),
	}
	err := suite.session.Insert(entity.StorageCollectionName, suite.storage)
	suite.NoError(err)
}

func (suite *VolumeTestSuite) TearDownSuite() {
	suite.session.Remove(entity.StorageCollectionName, "_id", suite.storage.ID)
}

func TestVolumeSuite(t *testing.T) {
	suite.Run(t, new(VolumeTestSuite))
}

func (suite *VolumeTestSuite) TestCreateVolume() {
	tName := namesgenerator.GetRandomName(0)
	tAccessMode := corev1.PersistentVolumeAccessMode("ReadOnlyMany")
	tCapacity := "500G"
	volume := entity.Volume{
		OwnerID:     bson.NewObjectId(),
		Name:        tName,
		StorageName: suite.storage.Name,
		Capacity:    tCapacity,
		AccessMode:  tAccessMode,
	}

	bodyBytes, err := json.MarshalIndent(volume, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/volume", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusCreated, httpWriter)
	defer suite.session.Remove(entity.VolumeCollectionName, "name", volume.Name)

	//load data to check
	retVolume := entity.Volume{}
	err = suite.session.FindOne(entity.VolumeCollectionName, bson.M{"name": volume.Name}, &retVolume)
	suite.NoError(err)
	suite.NotEqual("", retVolume.ID)
	suite.Equal(volume.Name, retVolume.Name)
	suite.Equal(volume.StorageName, retVolume.StorageName)
	suite.Equal(volume.AccessMode, retVolume.AccessMode)
	suite.Equal(volume.Capacity, retVolume.Capacity)

	//We use the new write but empty input which will cause the readEntity Error
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
	//Create again and it should fail since the name exist
	bodyReader = strings.NewReader(string(bodyBytes))
	httpRequest, err = http.NewRequest("POST", "http://localhost:7890/v1/volume", bodyReader)
	suite.NoError(err)
	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusConflict, httpWriter)
}

func (suite *VolumeTestSuite) TestCreateVolumeWithInvalidParameter() {
	//the storageName doesn't exist
	tName := namesgenerator.GetRandomName(0)
	tAccessMode := corev1.PersistentVolumeAccessMode("ReadOnlyMany")
	tCapacity := "500G"
	volume := entity.Volume{
		OwnerID:     bson.NewObjectId(),
		Name:        tName,
		StorageName: namesgenerator.GetRandomName(0),
		Capacity:    tCapacity,
		AccessMode:  tAccessMode,
	}

	bodyBytes, err := json.MarshalIndent(volume, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/volume", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusInternalServerError, httpWriter)
	defer suite.session.Remove(entity.VolumeCollectionName, "name", volume.Name)
}

func (suite *VolumeTestSuite) TestDeleteVolume() {
	tName := namesgenerator.GetRandomName(0)
	tAccessMode := corev1.PersistentVolumeAccessMode("ReadOnlyMany")
	tCapacity := "250"
	volume := entity.Volume{
		ID:          bson.NewObjectId(),
		OwnerID:     bson.NewObjectId(),
		Name:        tName,
		StorageName: namesgenerator.GetRandomName(0),
		Capacity:    tCapacity,
		AccessMode:  tAccessMode,
	}

	err := suite.session.Insert(entity.StorageCollectionName, &entity.Storage{
		Name: volume.StorageName,
	})
	defer suite.session.Remove(entity.StorageCollectionName, "name", volume.StorageName)
	suite.NoError(err)

	err = v.CreateVolume(suite.sp, &volume)
	suite.NoError(err)

	err = suite.session.Insert(entity.VolumeCollectionName, &volume)
	suite.NoError(err)

	bodyBytes, err := json.MarshalIndent(volume, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/volume/"+volume.ID.Hex(), bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	n, err := suite.session.Count(entity.VolumeCollectionName, bson.M{"_id": volume.ID})
	suite.NoError(err)
	suite.Equal(0, n)
}

func (suite *VolumeTestSuite) TestDeleteVolumeWithInvalidID() {
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/volume/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
}

func (suite *VolumeTestSuite) TestListVolume() {
	volumes := []entity.Volume{}
	count := 3
	for i := 0; i < count; i++ {
		volumes = append(volumes, entity.Volume{
			ID:          bson.NewObjectId(),
			OwnerID:     bson.NewObjectId(),
			Name:        namesgenerator.GetRandomName(0),
			StorageName: namesgenerator.GetRandomName(0),
			AccessMode:  corev1.PersistentVolumeAccessMode("ReadOnlyMany"),
			Capacity:    "250",
		})
	}

	for _, v := range volumes {
		suite.session.C(entity.VolumeCollectionName).Insert(v)
		defer suite.session.Remove(entity.VolumeCollectionName, "_id", v.ID)
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

	for _, tc := range testCases {
		caseName := "page:pageSize" + tc.page + ":" + tc.pageSize
		suite.T().Run(caseName, func(t *testing.T) {
			//list data by default page and page_size
			url := "http://localhost:7890/v1/volume/"
			if tc.page != "" || tc.pageSize != "" {
				url = "http://localhost:7890/v1/volume?"
				url += "page=" + tc.page + "%" + "page_size" + tc.pageSize
			}
			httpRequest, err := http.NewRequest("GET", url, nil)
			suite.NoError(err)

			httpRequest.Header.Add("Authorization", suite.JWTBearer)
			httpWriter := httptest.NewRecorder()
			suite.wc.Dispatch(httpWriter, httpRequest)
			assertResponseCode(suite.T(), http.StatusOK, httpWriter)

			retVolumes := []entity.Volume{}
			err = json.Unmarshal(httpWriter.Body.Bytes(), &retVolumes)
			suite.NoError(err)
			suite.Equal(tc.expectSize, len(retVolumes))
			for i, v := range retVolumes {
				suite.Equal(volumes[i].Name, v.Name)
				suite.Equal(volumes[i].StorageName, v.StorageName)
				suite.Equal(volumes[i].AccessMode, v.AccessMode)
			}
		})
	}
}

func (suite *VolumeTestSuite) TestListVolumeWithInvalidPage() {
	//Get data with non-exits ID
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/volume?page=asdd", nil)
	suite.NoError(err)

	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/volume?page_size=asdd", nil)
	suite.NoError(err)

	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/volume?page=-1", nil)
	suite.NoError(err)

	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusInternalServerError, httpWriter)
}
