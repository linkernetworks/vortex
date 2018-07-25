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

type PrometheusTestSuite struct {
	suite.Suite
	wc *restful.Container
	sp *serviceprovider.Container
}

func (suite *PrometheusTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	//init restful container
	suite.wc = restful.NewContainer()
	suite.sp = serviceprovider.New(cf)
	service := newMonitoringService(sp)
	suite.wc.Add(service)
}

func (suite *PrometheusTestSuite) TearDownSuite() {}

func TestPrometheusSuite(t *testing.T) {
	suite.Run(t, new(PrometheusTestSuite))
}

func (suite *PrometheusTestSuite) TestListNodeMetricsStatus() {
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/nodes/", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestGetNodeMetricsStatus() {
	nodes, err := suite.sp.KubeCtl.GetNodes()
	suite.NoError(err)
	nodeName := nodes[0].GetName()

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/nodes/"+nodeName, nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestListNodeNicsMetricsStatus() {
	nodes, err := suite.sp.KubeCtl.GetNodes()
	suite.NoError(err)
	nodeName := nodes[0].GetName()

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/nodes/"+nodeName+"/nics", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestListPodMetricsStatus() {
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/pods/", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestGetPodMetricsStatus() {
	namespace := "vortex"
	pods, err := suite.sp.KubeCtl.GetPods(namespace)
	suite.NoError(err)
	podName := pods[0].GetName()

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/pods/"+podName, nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestListContainerMetricsStatus() {
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/containers/", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestGetContainerMetricsStatus() {
	namespace := "vortex"
	pods, err := suite.sp.KubeCtl.GetPods(namespace)
	suite.NoError(err)
	containerName := pods[0].Status.ContainerStatuses[0].Name

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/containers/"+containerName, nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestListServiceMetricsStatus() {
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/services", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestGetServiceMetricsStatus() {
	namespace := "vortex"
	services, err := suite.sp.KubeCtl.GetServices(namespace)
	suite.NoError(err)
	serviceName := services[0].GetName()

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/services/"+serviceName, nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestListControllerMetricsStatus() {
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/controllers/", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestGetControllerMetricsStatus() {
	namespace := "vortex"
	deployments, err := suite.sp.KubeCtl.GetDeployments(namespace)
	suite.NoError(err)
	deploymentName := deployments[0].GetName()

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/controllers/"+deploymentName, nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}