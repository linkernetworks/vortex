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
	sp      *serviceprovider.Container
	wc      *restful.Container
	session *mongo.Session
}

func (suite *DeploymentTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	suite.sp = sp
	//init session
	suite.session = sp.Mongo.NewSession()
	//init restful container
	suite.wc = restful.NewContainer()
	service := newDeploymentService(suite.sp)
	suite.wc.Add(service)
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
	pod := entity.Deployment{
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
	bodyBytes, err := json.MarshalIndent(pod, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/deployments", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusCreated, httpWriter)
	defer suite.session.Remove(entity.DeploymentCollectionName, "name", pod.Name)

	//load data to check
	retDeployment := entity.Deployment{}
	err = suite.session.FindOne(entity.DeploymentCollectionName, bson.M{"name": pod.Name}, &retDeployment)
	suite.NoError(err)
	suite.NotEqual("", retDeployment.ID)
	suite.Equal(pod.Name, retDeployment.Name)
	suite.Equal(len(pod.Containers), len(retDeployment.Containers))

	//We use the new write but empty input which will cause the readEntity Error
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
	//Create again and it should fail since the name exist
	bodyReader = strings.NewReader(string(bodyBytes))
	httpRequest, err = http.NewRequest("POST", "http://localhost:7890/v1/deployments", bodyReader)
	suite.NoError(err)
	httpRequest.Header.Add("Content-Type", "application/json")
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
	pod := entity.Deployment{
		Name:       tName,
		Namespace:  namespace,
		Containers: containers,
		Volumes: []entity.DeploymentVolume{
			{Name: namesgenerator.GetRandomName(0)},
		},
	}

	bodyBytes, err := json.MarshalIndent(pod, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/deployments", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
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
	pod := entity.Deployment{
		ID:           bson.NewObjectId(),
		Name:         tName,
		Namespace:    namespace,
		Containers:   containers,
		Capability:   true,
		NetworkType:  entity.DeploymentHostNetwork,
		NodeAffinity: []string{},
		Replicas:     1,
	}

	err := p.CreateDeployment(suite.sp, &pod)
	suite.NoError(err)

	err = suite.session.Insert(entity.DeploymentCollectionName, &pod)
	suite.NoError(err)

	bodyBytes, err := json.MarshalIndent(pod, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/deployments/"+pod.ID.Hex(), bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	n, err := suite.session.Count(entity.DeploymentCollectionName, bson.M{"_id": pod.ID})
	suite.NoError(err)
	suite.Equal(0, n)
}

func (suite *DeploymentTestSuite) TestDeleteDeploymentWithInvalidID() {
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/deployments/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
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
	pod := entity.Deployment{
		ID:         bson.NewObjectId(),
		Name:       tName,
		Namespace:  namespace,
		Containers: containers,
	}

	//Create data into mongo manually
	suite.session.C(entity.DeploymentCollectionName).Insert(pod)
	defer suite.session.Remove(entity.DeploymentCollectionName, "name", tName)

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/deployments/"+pod.ID.Hex(), nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	pod = entity.Deployment{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &pod)
	suite.NoError(err)
	suite.Equal(tName, pod.Name)
	suite.Equal(len(containers), len(pod.Containers))
}

func (suite *DeploymentTestSuite) TestGetDeploymentWithInvalidID() {
	//Get data with non-exits ID
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/deployments/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)

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

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/deployments?page_size=asdd", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/deployments?page=-1", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusInternalServerError, httpWriter)
}
