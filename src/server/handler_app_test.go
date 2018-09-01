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

type AppTestSuite struct {
	suite.Suite
	sp        *serviceprovider.Container
	wc        *restful.Container
	session   *mongo.Session
	JWTBearer string
}

func (suite *AppTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	suite.sp = sp
	//init session
	suite.session = sp.Mongo.NewSession()
	//init restful container
	suite.wc = restful.NewContainer()
	appService := newAppService(suite.sp)
	userService := newUserService(suite.sp)

	suite.wc.Add(appService)
	suite.wc.Add(userService)

	token, _ := loginGetToken(suite.wc)
	suite.NotEmpty(token)
	suite.JWTBearer = "Bearer " + token
}

func (suite *AppTestSuite) TearDownSuite() {}

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}

func (suite *AppTestSuite) TestCreateApp() {
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

	ports := []entity.ServicePort{
		{
			Name:       namesgenerator.GetRandomName(0),
			Port:       int32(80),
			TargetPort: 80,
			NodePort:   int32(30000),
		},
	}

	serviceName := namesgenerator.GetRandomName(0)
	service := entity.Service{
		ID:        bson.NewObjectId(),
		Name:      serviceName,
		Namespace: "default",
		Type:      "NodePort",
		Selector:  map[string]string{},
		Ports:     ports,
	}

	app := entity.Application{
		Deployment: deploy,
		Service:    service,
	}
	bodyBytes, err := json.MarshalIndent(app, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/apps", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusCreated, httpWriter)
	defer suite.session.Remove(entity.DeploymentCollectionName, "name", app.Deployment.Name)
	defer suite.session.Remove(entity.ServiceCollectionName, "name", app.Service.Name)

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
	httpRequest, err = http.NewRequest("POST", "http://localhost:7890/v1/apps", bodyReader)
	suite.NoError(err)
	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusConflict, httpWriter)
}
