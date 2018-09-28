package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/configmap"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type ConfigMapTestSuite struct {
	suite.Suite
	sp        *serviceprovider.Container
	wc        *restful.Container
	session   *mongo.Session
	JWTBearer string
}

func (suite *ConfigMapTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	suite.sp = sp
	// init session
	suite.session = sp.Mongo.NewSession()
	// init restful container
	suite.wc = restful.NewContainer()

	configMapService := newConfigMapService(suite.sp)
	userService := newUserService(suite.sp)

	suite.wc.Add(configMapService)
	suite.wc.Add(userService)

	token, _ := loginGetToken(suite.wc)
	suite.NotEmpty(token)
	suite.JWTBearer = "Bearer " + token
}

func (suite *ConfigMapTestSuite) TearDownSuite() {}

func TestConfigMapSuite(t *testing.T) {
	suite.Run(t, new(ConfigMapTestSuite))
}

func (suite *ConfigMapTestSuite) TestCreateConfigMap() {
	data := map[string]string{
		"firstData":  "YXdlc29tZQ==",
		"secondData": "ewogICJjb2xvcnMiOiBbCiAgICB7CiAgICAgICJjb2xvciI6ICJibGFjayIsCiAgICAgICJjYXRlZ29yeSI6ICJodWUiLAogICAgICAidHlwZSI6ICJwcmltYXJ5IiwKICAgICAgImNvZGUiOiB7CiAgICAgICAgInJnYmEiOiBbMjU1LDI1NSwyNTUsMV0sCiAgICAgICAgImhleCI6ICIjMDAwIgogICAgICB9CiAgICB9LAogICAgewogICAgICAiY29sb3IiOiAid2hpdGUiLAogICAgICAiY2F0ZWdvcnkiOiAidmFsdWUiLAogICAgICAiY29kZSI6IHsKICAgICAgICAicmdiYSI6IFswLDAsMCwxXSwKICAgICAgICAiaGV4IjogIiNGRkYiCiAgICAgIH0KICAgIH0sCiAgICB7CiAgICAgICJjb2xvciI6ICJyZWQiLAogICAgICAiY2F0ZWdvcnkiOiAiaHVlIiwKICAgICAgInR5cGUiOiAicHJpbWFyeSIsCiAgICAgICJjb2RlIjogewogICAgICAgICJyZ2JhIjogWzI1NSwwLDAsMV0sCiAgICAgICAgImhleCI6ICIjRkYwIgogICAgICB9CiAgICB9LAogICAgewogICAgICAiY29sb3IiOiAiYmx1ZSIsCiAgICAgICJjYXRlZ29yeSI6ICJodWUiLAogICAgICAidHlwZSI6ICJwcmltYXJ5IiwKICAgICAgImNvZGUiOiB7CiAgICAgICAgInJnYmEiOiBbMCwwLDI1NSwxXSwKICAgICAgICAiaGV4IjogIiMwMEYiCiAgICAgIH0KICAgIH0sCiAgICB7CiAgICAgICJjb2xvciI6ICJ5ZWxsb3ciLAogICAgICAiY2F0ZWdvcnkiOiAiaHVlIiwKICAgICAgInR5cGUiOiAicHJpbWFyeSIsCiAgICAgICJjb2RlIjogewogICAgICAgICJyZ2JhIjogWzI1NSwyNTUsMCwxXSwKICAgICAgICAiaGV4IjogIiNGRjAiCiAgICAgIH0KICAgIH0sCiAgICB7CiAgICAgICJjb2xvciI6ICJncmVlbiIsCiAgICAgICJjYXRlZ29yeSI6ICJodWUiLAogICAgICAidHlwZSI6ICJzZWNvbmRhcnkiLAogICAgICAiY29kZSI6IHsKICAgICAgICAicmdiYSI6IFswLDI1NSwwLDFdLAogICAgICAgICJoZXgiOiAiIzBGMCIKICAgICAgfQogICAgfSwKICBdCn0=",
	}
	configMap := entity.ConfigMap{
		ID:        bson.NewObjectId(),
		Name:      namesgenerator.GetRandomName(0),
		Namespace: "default",
		Data:      data,
	}

	bodyBytes, err := json.MarshalIndent(configMap, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/configMaps", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusCreated, httpWriter)
	defer suite.session.Remove(entity.ConfigMapCollectionName, "name", configMap.Name)

	//load data to check
	retConfigMap := entity.ConfigMap{}
	err = suite.session.FindOne(entity.ConfigMapCollectionName, bson.M{"name": configMap.Name}, &retConfigMap)
	suite.NoError(err)
	suite.NotEqual("", retConfigMap.ID)
	suite.Equal(configMap.Name, retConfigMap.Name)

	//We use the new write but empty input which will cause the readEntity Error
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	//Create again and it should fail since the name exist
	bodyReader = strings.NewReader(string(bodyBytes))
	httpRequest, err = http.NewRequest("POST", "http://localhost:7890/v1/configMaps", bodyReader)
	suite.NoError(err)
	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusConflict, httpWriter)
	defer suite.session.Remove(entity.ConfigMapCollectionName, "name", configMap.Name)

	err = configmap.DeleteConfigMap(suite.sp, &configMap)
	suite.NoError(err)
}

func (suite *ConfigMapTestSuite) TestDeleteConfigMap() {
	configMapName := namesgenerator.GetRandomName(0)
	configMap := entity.ConfigMap{
		ID:   bson.NewObjectId(),
		Name: configMapName,
	}

	err := configmap.CreateConfigMap(suite.sp, &configMap)
	suite.NoError(err)

	err = suite.session.Insert(entity.ConfigMapCollectionName, &configMap)
	suite.NoError(err)

	bodyBytes, err := json.MarshalIndent(configMap, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/configMaps/"+configMap.ID.Hex(), bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	n, err := suite.session.Count(entity.ConfigMapCollectionName, bson.M{"_id": configMap.ID})
	suite.NoError(err)
	suite.Equal(0, n)
}

func (suite *ConfigMapTestSuite) TestDeleteConfigMapWithInvalidID() {
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/configMaps/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
}

//For Get/List, we only return mongo document
func (suite *ConfigMapTestSuite) TestGetConfigMap() {
	configMapName := namesgenerator.GetRandomName(0)
	data := map[string]string{
		"firstData":  "YXdlc29tZQ==",
		"secondData": "ewogICJjb2xvcnMiOiBbCiAgICB7CiAgICAgICJjb2xvciI6ICJibGFjayIsCiAgICAgICJjYXRlZ29yeSI6ICJodWUiLAogICAgICAidHlwZSI6ICJwcmltYXJ5IiwKICAgICAgImNvZGUiOiB7CiAgICAgICAgInJnYmEiOiBbMjU1LDI1NSwyNTUsMV0sCiAgICAgICAgImhleCI6ICIjMDAwIgogICAgICB9CiAgICB9LAogICAgewogICAgICAiY29sb3IiOiAid2hpdGUiLAogICAgICAiY2F0ZWdvcnkiOiAidmFsdWUiLAogICAgICAiY29kZSI6IHsKICAgICAgICAicmdiYSI6IFswLDAsMCwxXSwKICAgICAgICAiaGV4IjogIiNGRkYiCiAgICAgIH0KICAgIH0sCiAgICB7CiAgICAgICJjb2xvciI6ICJyZWQiLAogICAgICAiY2F0ZWdvcnkiOiAiaHVlIiwKICAgICAgInR5cGUiOiAicHJpbWFyeSIsCiAgICAgICJjb2RlIjogewogICAgICAgICJyZ2JhIjogWzI1NSwwLDAsMV0sCiAgICAgICAgImhleCI6ICIjRkYwIgogICAgICB9CiAgICB9LAogICAgewogICAgICAiY29sb3IiOiAiYmx1ZSIsCiAgICAgICJjYXRlZ29yeSI6ICJodWUiLAogICAgICAidHlwZSI6ICJwcmltYXJ5IiwKICAgICAgImNvZGUiOiB7CiAgICAgICAgInJnYmEiOiBbMCwwLDI1NSwxXSwKICAgICAgICAiaGV4IjogIiMwMEYiCiAgICAgIH0KICAgIH0sCiAgICB7CiAgICAgICJjb2xvciI6ICJ5ZWxsb3ciLAogICAgICAiY2F0ZWdvcnkiOiAiaHVlIiwKICAgICAgInR5cGUiOiAicHJpbWFyeSIsCiAgICAgICJjb2RlIjogewogICAgICAgICJyZ2JhIjogWzI1NSwyNTUsMCwxXSwKICAgICAgICAiaGV4IjogIiNGRjAiCiAgICAgIH0KICAgIH0sCiAgICB7CiAgICAgICJjb2xvciI6ICJncmVlbiIsCiAgICAgICJjYXRlZ29yeSI6ICJodWUiLAogICAgICAidHlwZSI6ICJzZWNvbmRhcnkiLAogICAgICAiY29kZSI6IHsKICAgICAgICAicmdiYSI6IFswLDI1NSwwLDFdLAogICAgICAgICJoZXgiOiAiIzBGMCIKICAgICAgfQogICAgfSwKICBdCn0=",
	}
	configMap := entity.ConfigMap{
		ID:        bson.NewObjectId(),
		Name:      configMapName,
		Namespace: "default",
		Data:      data,
	}

	// Create data into mongo manually
	suite.session.C(entity.ConfigMapCollectionName).Insert(configMap)
	defer suite.session.Remove(entity.ConfigMapCollectionName, "name", configMapName)

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/configMaps/"+configMap.ID.Hex(), nil)
	suite.NoError(err)

	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	configMap = entity.ConfigMap{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &configMap)
	suite.NoError(err)
	suite.Equal(configMapName, configMap.Name)
}

func (suite *ConfigMapTestSuite) TestGetConfigMapWithInvalidID() {
	// Get data with non-exits ID
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/configMaps/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)

	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusNotFound, httpWriter)
}

func (suite *ConfigMapTestSuite) TestListConfigMap() {
	configMaps := []entity.ConfigMap{}
	count := 3
	for i := 0; i < count; i++ {
		configMaps = append(configMaps, entity.ConfigMap{
			ID:        bson.NewObjectId(),
			Name:      namesgenerator.GetRandomName(0),
			Namespace: "default",
		})
	}

	for _, n := range configMaps {
		suite.session.C(entity.ConfigMapCollectionName).Insert(n)
		defer suite.session.Remove(entity.ConfigMapCollectionName, "_id", n.ID)
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
			url := "http://localhost:7890/v1/configMaps/"
			if tc.page != "" || tc.pageSize != "" {
				url = "http://localhost:7890/v1/configMaps?"
				url += "page=" + tc.page + "%" + "page_size" + tc.pageSize
			}
			httpRequest, err := http.NewRequest("GET", url, nil)
			suite.NoError(err)

			httpRequest.Header.Add("Authorization", suite.JWTBearer)
			httpWriter := httptest.NewRecorder()
			suite.wc.Dispatch(httpWriter, httpRequest)
			assertResponseCode(suite.T(), http.StatusOK, httpWriter)

			retConfigMaps := []entity.ConfigMap{}
			err = json.Unmarshal(httpWriter.Body.Bytes(), &retConfigMaps)
			suite.NoError(err)
			suite.Equal(tc.expectSize, len(retConfigMaps))
			for i, n := range retConfigMaps {
				suite.Equal(configMaps[i].Name, n.Name)
			}
		})
	}
}

func (suite *ConfigMapTestSuite) TestListConfigMapWithInvalidPage() {
	//Get data with non-exits ID
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/configMaps?page=asdd", nil)
	suite.NoError(err)

	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/configMaps?page_size=asdd", nil)
	suite.NoError(err)

	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/configMaps?page=-1", nil)
	suite.NoError(err)

	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusInternalServerError, httpWriter)
}

func (suite *ConfigMapTestSuite) TestUploadConfigMapYAML() {
	filename := "../../testYAMLs/configMap.yaml"

	bodyBuf := bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(bodyBuf)

	// use the bodyWriter to write the Part headers to the buffer
	_, err := bodyWriter.CreateFormFile("file", filename)
	suite.NoError(err)

	// the file data will be the second part of the body
	file, err := os.Open(filename)
	suite.NoError(err)

	// need to know the boundary to properly close the part myself.
	boundary := bodyWriter.Boundary()
	//close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until
	// writing to the socket buffer.
	requestReader := io.MultiReader(bodyBuf, file, closeBuf)
	fileStat, err := file.Stat()
	suite.NoError(err)

	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/configMaps/upload/yaml", requestReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpRequest.ContentLength = fileStat.Size() + int64(bodyBuf.Len()) + int64(closeBuf.Len())
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	defer suite.session.Remove(entity.ConfigMapCollectionName, "name", "upload-configmap")

	assertResponseCode(suite.T(), http.StatusCreated, httpWriter)

	//load data to check
	retConfigMap := entity.ConfigMap{}
	err = suite.session.FindOne(entity.ConfigMapCollectionName, bson.M{"name": "upload-configmap"}, &retConfigMap)
	suite.NoError(err)
	suite.NotEqual("", retConfigMap.ID)
	suite.Equal("upload-configmap", retConfigMap.Name)

	//Create again and it should fail since the name exist
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusConflict, httpWriter)

	err = configmap.DeleteConfigMap(suite.sp, &retConfigMap)
	suite.NoError(err)
}

func (suite *ConfigMapTestSuite) TestUploadConfigMapYAMLFail() {
	filename := "../../testYAMLs/deployment.yaml"

	bodyBuf := bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(bodyBuf)

	// use the bodyWriter to write the Part headers to the buffer
	_, err := bodyWriter.CreateFormFile("file", filename)
	suite.NoError(err)

	// the file data will be the second part of the body
	file, err := os.Open(filename)
	suite.NoError(err)

	// need to know the boundary to properly close the part myself.
	boundary := bodyWriter.Boundary()
	//close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until
	// writing to the socket buffer.
	requestReader := io.MultiReader(bodyBuf, file, closeBuf)
	fileStat, err := file.Stat()
	suite.NoError(err)

	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/configMaps/upload/yaml", requestReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpRequest.ContentLength = fileStat.Size() + int64(bodyBuf.Len()) + int64(closeBuf.Len())
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)

	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
}
