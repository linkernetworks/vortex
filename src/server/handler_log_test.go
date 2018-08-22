package server

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/stretchr/testify/suite"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type LogTestSuite struct {
	suite.Suite
	sp *serviceprovider.Container
	wc *restful.Container
}

func (suite *LogTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	suite.sp = sp
	// init restful container
	suite.wc = restful.NewContainer()
	service := newContainerService(suite.sp)
	suite.wc.Add(service)
}

func (suite *LogTestSuite) TearDownSuite() {}

func TestLogSuite(t *testing.T) {
	suite.Run(t, new(LogTestSuite))
}

func (suite *LogTestSuite) TestGetContainerLogs() {
	namespace := "vortex"
	pods, err := suite.sp.KubeCtl.GetPods(namespace)
	suite.NoError(err)
	podName := pods[0].Name
	containerName := pods[0].Status.ContainerStatuses[0].Name

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/containers/logs/"+namespace+"/"+podName+"/"+containerName, nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *LogTestSuite) TestGetContainerLogFile() {
	namespace := "vortex"
	pods, err := suite.sp.KubeCtl.GetPods(namespace)
	suite.NoError(err)
	podName := pods[0].Name
	containerName := pods[0].Status.ContainerStatuses[0].Name

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/containers/logs/file/"+namespace+"/"+podName+"/"+containerName, nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}
