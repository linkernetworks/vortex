package prometheuscontroller

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/suite"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type PrometheusExpressionTestSuite struct {
	suite.Suite
	sp *serviceprovider.Container
}

func (suite *PrometheusExpressionTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/local.json")
	suite.sp = serviceprovider.New(cf)
	suite.sp.KubeCtl.SetNamespace("monitoring")
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

func (suite *PrometheusExpressionTestSuite) TestListResource() {
	labelName := model.LabelName("container")
	expression := Expression{}
	expression.Metrics = []string{"kube_pod_container_info"}
	expression.QueryLabels = map[string]string{}
	expression.QueryLabels["namespace"] = "monitoring"

	resourceList, err := ListResource(suite.sp, labelName, expression)
	suite.NoError(err)
	suite.NotEqual(0, len(resourceList))
}

func (suite *PrometheusExpressionTestSuite) TestListResourceFail() {
	labelName := model.LabelName("")
	expression := Expression{}

	_, err := ListResource(suite.sp, labelName, expression)
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
	namespace := "default"
	pods, err := suite.sp.KubeCtl.GetPods(namespace)
	suite.NoError(err)
	podName := pods[0].GetName()

	pod, err := GetPod(suite.sp, podName)
	suite.NoError(err)
	suite.Equal(podName, pod.PodName)
}

func (suite *PrometheusExpressionTestSuite) TestGetContainer() {
	namespace := "default"
	pods, err := suite.sp.KubeCtl.GetPods(namespace)
	suite.NoError(err)
	containerName := pods[0].Status.ContainerStatuses[0].Name

	container, err := GetContainer(suite.sp, containerName)
	suite.NoError(err)
	suite.Equal(containerName, container.Detail.ContainerName)
}

func (suite *PrometheusExpressionTestSuite) TestGetService() {
	namespace := "default"
	services, err := suite.sp.KubeCtl.GetServices(namespace)
	suite.NoError(err)
	serviceName := services[0].GetName()

	service, err := GetService(suite.sp, serviceName)
	suite.NoError(err)
	suite.Equal(serviceName, service.ServiceName)
}

func (suite *PrometheusExpressionTestSuite) TestGetController() {
	namespace := "default"
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

	node, err := GetNode(suite.sp, nodeName)
	suite.NoError(err)
	suite.Equal(nodeName, node.Detail.Hostname)
}
