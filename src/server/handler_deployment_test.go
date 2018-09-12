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
	p "github.com/linkernetworks/vortex/src/deployment"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type DeploymentTestSuite struct {
	suite.Suite
	sp        *serviceprovider.Container
	wc        *restful.Container
	session   *mongo.Session
	JWTBearer string
}

func (suite *DeploymentTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	suite.sp = sp
	// init session
	suite.session = sp.Mongo.NewSession()
	// init restful container
	suite.wc = restful.NewContainer()

	deploymentService := newDeploymentService(suite.sp)
	userService := newUserService(suite.sp)

	suite.wc.Add(deploymentService)
	suite.wc.Add(userService)

	token, _ := loginGetToken(suite.wc)
	suite.NotEmpty(token)
	suite.JWTBearer = "Bearer " + token
}

func (suite *DeploymentTestSuite) TearDownSuite() {}

func TestDeploymentSuite(t *testing.T) {
	suite.Run(t, new(DeploymentTestSuite))
}

func (suite *DeploymentTestSuite) TestCreateDeployment() {
	namespace := "default"
	containers := []entity.Container{
		{
			Name:    namesgenerator.GetRandomName(0),
			Image:   "busybox",
			Command: []string{"sleep", "3600"},
		},
	}
	tName := namesgenerator.GetRandomName(0)
	deploy := entity.Deployment{
		Name:         tName,
		Namespace:    namespace,
		Labels:       map[string]string{},
		EnvVars:      map[string]string{},
		Containers:   containers,
		Volumes:      []entity.DeploymentVolume{},
		Networks:     []entity.DeploymentNetwork{},
		Capability:   true,
		NetworkType:  entity.DeploymentHostNetwork,
		NodeAffinity: []string{},
		Replicas:     1,
	}
	bodyBytes, err := json.MarshalIndent(deploy, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/deployments", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusCreated, httpWriter)
	defer suite.session.Remove(entity.DeploymentCollectionName, "name", deploy.Name)

	//load data to check
	retDeployment := entity.Deployment{}
	err = suite.session.FindOne(entity.DeploymentCollectionName, bson.M{"name": deploy.Name}, &retDeployment)
	suite.NoError(err)
	suite.NotEqual("", retDeployment.ID)
	suite.Equal(deploy.Name, retDeployment.Name)
	suite.Equal(len(deploy.Containers), len(retDeployment.Containers))

	//We use the new write but empty input which will cause the readEntity Error
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
	//Create again and it should fail since the name exist
	bodyReader = strings.NewReader(string(bodyBytes))
	httpRequest, err = http.NewRequest("POST", "http://localhost:7890/v1/deployments", bodyReader)
	suite.NoError(err)
	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusConflict, httpWriter)

	err = p.DeleteDeployment(suite.sp, &retDeployment)
	suite.NoError(err)
}

func (suite *DeploymentTestSuite) TestCreateDeploymentFail() {
	namespace := "default"
	containers := []entity.Container{
		{
			Name:    namesgenerator.GetRandomName(0),
			Image:   "busybox",
			Command: []string{"sleep", "3600"},
		},
	}
	tName := namesgenerator.GetRandomName(0)
	deploy := entity.Deployment{
		Name:       tName,
		Namespace:  namespace,
		Containers: containers,
		Volumes: []entity.DeploymentVolume{
			{Name: namesgenerator.GetRandomName(0)},
		},
	}

	bodyBytes, err := json.MarshalIndent(deploy, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/deployments", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
}

func (suite *DeploymentTestSuite) TestDeleteDeployment() {
	namespace := "default"
	containers := []entity.Container{
		{
			Name:    namesgenerator.GetRandomName(0),
			Image:   "busybox",
			Command: []string{"sleep", "3600"},
		},
	}
	tName := namesgenerator.GetRandomName(0)
	deploy := entity.Deployment{
		ID:           bson.NewObjectId(),
		Name:         tName,
		Namespace:    namespace,
		Containers:   containers,
		Capability:   true,
		NetworkType:  entity.DeploymentHostNetwork,
		NodeAffinity: []string{},
		Replicas:     1,
	}

	err := p.CreateDeployment(suite.sp, &deploy)
	suite.NoError(err)

	err = suite.session.Insert(entity.DeploymentCollectionName, &deploy)
	suite.NoError(err)

	bodyBytes, err := json.MarshalIndent(deploy, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/deployments/"+deploy.ID.Hex(), bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	n, err := suite.session.Count(entity.DeploymentCollectionName, bson.M{"_id": deploy.ID})
	suite.NoError(err)
	suite.Equal(0, n)
}

func (suite *DeploymentTestSuite) TestDeleteDeploymentWithInvalidID() {
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/deployments/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
}

//For Get/List, we only return mongo document
func (suite *DeploymentTestSuite) TestGetDeployment() {
	namespace := "default"
	containers := []entity.Container{
		{
			Name:    namesgenerator.GetRandomName(0),
			Image:   "busybox",
			Command: []string{"sleep", "3600"},
		},
	}
	tName := namesgenerator.GetRandomName(0)
	deploy := entity.Deployment{
		ID:         bson.NewObjectId(),
		Name:       tName,
		Namespace:  namespace,
		Containers: containers,
	}

	//Create data into mongo manually
	suite.session.C(entity.DeploymentCollectionName).Insert(deploy)
	defer suite.session.Remove(entity.DeploymentCollectionName, "name", tName)

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/deployments/"+deploy.ID.Hex(), nil)
	suite.NoError(err)

	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	deploy = entity.Deployment{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &deploy)
	suite.NoError(err)
	suite.Equal(tName, deploy.Name)
	suite.Equal(len(containers), len(deploy.Containers))
}

func (suite *DeploymentTestSuite) TestGetDeploymentWithInvalidID() {
	//Get data with non-exits ID
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/deployments/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)

	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusNotFound, httpWriter)
}

func (suite *DeploymentTestSuite) TestListDeployment() {
	namespace := "default"
	deployments := []entity.Deployment{}
	count := 3
	for i := 0; i < count; i++ {
		containers := []entity.Container{
			{
				Name:    namesgenerator.GetRandomName(0),
				Image:   "busybox",
				Command: []string{"sleep", "3600"},
			},
		}
		deployments = append(deployments, entity.Deployment{
			ID:         bson.NewObjectId(),
			Name:       namesgenerator.GetRandomName(0),
			Namespace:  namespace,
			Containers: containers,
		})
	}

	for _, p := range deployments {
		suite.session.C(entity.DeploymentCollectionName).Insert(p)
		defer suite.session.Remove(entity.DeploymentCollectionName, "_id", p.ID)
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
			url := "http://localhost:7890/v1/deployments/"
			if tc.page != "" || tc.pageSize != "" {
				url = "http://localhost:7890/v1/deployments?"
				url += "page=" + tc.page + "%" + "page_size" + tc.pageSize
			}
			httpRequest, err := http.NewRequest("GET", url, nil)
			suite.NoError(err)

			httpRequest.Header.Add("Authorization", suite.JWTBearer)
			httpWriter := httptest.NewRecorder()
			suite.wc.Dispatch(httpWriter, httpRequest)
			assertResponseCode(suite.T(), http.StatusOK, httpWriter)

			retDeployments := []entity.Deployment{}
			err = json.Unmarshal(httpWriter.Body.Bytes(), &retDeployments)
			suite.NoError(err)
			suite.Equal(tc.expectSize, len(retDeployments))
			for i, p := range retDeployments {
				suite.Equal(deployments[i].Name, p.Name)
				suite.Equal(len(deployments[i].Containers), len(p.Containers))
			}
		})
	}
}

func (suite *DeploymentTestSuite) TestListDeploymentWithInvalidPage() {
	//Get data with non-exits ID
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/deployments?page=asdd", nil)
	suite.NoError(err)

	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/deployments?page_size=asdd", nil)
	suite.NoError(err)

	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/deployments?page=-1", nil)
	suite.NoError(err)

	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusInternalServerError, httpWriter)
}

func (suite *DeploymentTestSuite) TestUploadDeploymentYAML() {
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

	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/deployments/upload/yaml", requestReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpRequest.ContentLength = fileStat.Size() + int64(bodyBuf.Len()) + int64(closeBuf.Len())
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	defer suite.session.Remove(entity.DeploymentCollectionName, "name", "upload-deployment")

	assertResponseCode(suite.T(), http.StatusCreated, httpWriter)

	//load data to check
	retDeployment := entity.Deployment{}
	err = suite.session.FindOne(entity.DeploymentCollectionName, bson.M{"name": "upload-deployment"}, &retDeployment)
	suite.NoError(err)
	suite.NotEqual("", retDeployment.ID)
	suite.Equal("upload-deployment", retDeployment.Name)
	suite.Equal("default", retDeployment.Namespace)

	//Create again and it should fail since the name exist
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusConflict, httpWriter)

	err = p.DeleteDeployment(suite.sp, &retDeployment)
	suite.NoError(err)
}

func (suite *DeploymentTestSuite) TestUploadDeploymentYAMLFail() {
	filename := "../../testYAMLs/namespace.yaml"

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

	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/deployments/upload/yaml", requestReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpRequest.ContentLength = fileStat.Size() + int64(bodyBuf.Len()) + int64(closeBuf.Len())
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)

	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
}
