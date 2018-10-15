package prometheuscontroller

import (
	"os"
	"testing"

	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/stretchr/testify/suite"
)

type PrometheusExpressionTestSuite struct {
	suite.Suite
	sp *serviceprovider.Container
}

func (suite *PrometheusExpressionTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	suite.sp = serviceprovider.New(cf)
}

func (suite *PrometheusExpressionTestSuite) TearDownSuite() {
}

func TestPrometheusExpressionSuite(t *testing.T) {
	if _, defined := os.LookupEnv("TEST_PROMETHEUS"); !defined {
		t.SkipNow()
		return
	}
	suite.Run(t, new(PrometheusExpressionTestSuite))
}

func (suite *PrometheusExpressionTestSuite) TestListContainer() {
	queryLabels := map[string]string{"namespace": "vortex"}

	containerNameList, err := ListContainerName(suite.sp, queryLabels)
	suite.NoError(err)
	suite.NotEqual(0, len(containerNameList))
}

func (suite *PrometheusExpressionTestSuite) TestListContainerFail() {
	queryLabels := map[string]string{"": ""}

	_, err := ListContainerName(suite.sp, queryLabels)
	suite.Error(err)
}

func (suite *PrometheusExpressionTestSuite) TestListPod() {
	queryLabels := map[string]string{"namespace": "vortex"}

	podNameList, err := ListPodName(suite.sp, queryLabels)
	suite.NoError(err)
	suite.NotEqual(0, len(podNameList))
}

func (suite *PrometheusExpressionTestSuite) TestListPodFail() {
	queryLabels := map[string]string{"": ""}

	_, err := ListPodName(suite.sp, queryLabels)
	suite.Error(err)
}

func (suite *PrometheusExpressionTestSuite) TestListService() {
	queryLabels := map[string]string{"namespace": "vortex"}

	serviceNameList, err := ListServiceName(suite.sp, queryLabels)
	suite.NoError(err)
	suite.NotEqual(0, len(serviceNameList))
}

func (suite *PrometheusExpressionTestSuite) TestListServiceFail() {
	queryLabels := map[string]string{"": ""}

	_, err := ListServiceName(suite.sp, queryLabels)
	suite.Error(err)
}

func (suite *PrometheusExpressionTestSuite) TestListController() {
	queryLabels := map[string]string{"namespace": "vortex"}

	controllerNameList, err := ListControllerName(suite.sp, queryLabels)
	suite.NoError(err)
	suite.NotEqual(0, len(controllerNameList))
}

func (suite *PrometheusExpressionTestSuite) TestListControllerFail() {
	queryLabels := map[string]string{"": ""}

	_, err := ListControllerName(suite.sp, queryLabels)
	suite.Error(err)
}

func (suite *PrometheusExpressionTestSuite) TestListNode() {
	queryLabels := map[string]string{}

	nodeNameList, err := ListNodeName(suite.sp, queryLabels)
	suite.NoError(err)
	suite.NotEqual(0, len(nodeNameList))
}

func (suite *PrometheusExpressionTestSuite) TestListNodeFail() {
	queryLabels := map[string]string{"": ""}

	_, err := ListNodeName(suite.sp, queryLabels)
	suite.Error(err)
}

func (suite *PrometheusExpressionTestSuite) TestListNodeNICs() {
	nodes, err := suite.sp.KubeCtl.GetNodes()
	suite.NoError(err)
	nodeName := nodes[0].GetName()

	nicList, err := ListNodeNICs(suite.sp, nodeName)
	suite.NoError(err)
	suite.NotEqual(0, len(nicList.NICs))
}

func (suite *PrometheusExpressionTestSuite) TestGetPod() {
	namespace := "vortex"
	pods, err := suite.sp.KubeCtl.GetPods(namespace)
	suite.NoError(err)
	podName := pods[0].GetName()

	rs := RangeSetting{Interval: 1, Resolution: 1, Rate: 1}
	pod, err := GetPod(suite.sp, podName, rs)
	suite.NoError(err)
	suite.Equal(podName, pod.PodName)
}

func (suite *PrometheusExpressionTestSuite) TestGetContainer() {
	namespace := "vortex"
	pods, err := suite.sp.KubeCtl.GetPods(namespace)
	suite.NoError(err)
	podName := pods[0].Name
	containerName := pods[0].Status.ContainerStatuses[0].Name

	rs := RangeSetting{Interval: 1, Resolution: 1, Rate: 1}
	container, err := GetContainer(suite.sp, podName, containerName, rs)
	suite.NoError(err)
	suite.Equal(containerName, container.Detail.ContainerName)
}

func (suite *PrometheusExpressionTestSuite) TestGetService() {
	namespace := "vortex"
	services, err := suite.sp.KubeCtl.GetServices(namespace)
	suite.NoError(err)
	serviceName := services[0].GetName()

	service, err := GetService(suite.sp, serviceName)
	suite.NoError(err)
	suite.Equal(serviceName, service.ServiceName)
}

func (suite *PrometheusExpressionTestSuite) TestGetController() {
	namespace := "vortex"
	deployments, err := suite.sp.KubeCtl.GetDeployments(namespace)
	suite.NoError(err)
	deploymentName := deployments[0].GetName()

	controller, err := GetController(suite.sp, deploymentName)
	suite.NoError(err)
	suite.Equal(deploymentName, controller.ControllerName)
}

func (suite *PrometheusExpressionTestSuite) TestGetNode() {
	nodes, err := suite.sp.KubeCtl.GetNodes()
	suite.NoError(err)
	nodeName := nodes[0].GetName()

	rs := RangeSetting{Interval: 1, Resolution: 1, Rate: 1}
	node, err := GetNode(suite.sp, nodeName, rs)
	suite.NoError(err)
	suite.Equal(nodeName, node.Detail.Hostname)
}
