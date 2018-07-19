package prometheuscontroller

import (
	"fmt"
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

type PrometheusQueryTestSuite struct {
	suite.Suite
	sp            *serviceprovider.Container
	containerName string
}

func (suite *PrometheusQueryTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	suite.sp = serviceprovider.New(cf)
	suite.containerName = "cadvisor"
}

func (suite *PrometheusQueryTestSuite) TearDownSuite() {
}

func TestPrometheusQuerySuite(t *testing.T) {
	if _, defined := os.LookupEnv("TEST_PROMETHEUS"); !defined {
		t.SkipNow()
		return
	}
	suite.Run(t, new(PrometheusQueryTestSuite))
}

func (suite *PrometheusQueryTestSuite) TestQuery() {
	queryStr := fmt.Sprintf(`sum(rate(container_cpu_usage_seconds_total{container_label_io_kubernetes_container_name=~"%s"}[1m])) * 100`, suite.containerName)

	resource, err := query(suite.sp, queryStr)
	suite.NoError(err)
	suite.NotEqual(0, float32(resource[0].Value))

	// Get nil if the resource is empty
	resource, err = query(suite.sp, "")
	suite.Equal(model.Vector(nil), resource)
}

func (suite *PrometheusQueryTestSuite) TestGetElements() {
	expression := Expression{}
	expression.Metrics = []string{"kube_pod_container_info"}
	expression.QueryLabels = map[string]string{}
	expression.QueryLabels["namespace"] = "monitoring"

	resource, err := getElements(suite.sp, expression)
	suite.NoError(err)
	suite.NotEqual(0, float32(resource[0].Value))

	// Get nil if the resource is empty
	resource, err = getElements(suite.sp, Expression{})
	suite.Equal(model.Vector(nil), resource)
}
