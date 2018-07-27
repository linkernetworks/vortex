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

type PrometheusExprTestSuite struct {
	suite.Suite
	sp            *serviceprovider.Container
	containerName string
}

func (suite *PrometheusExprTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	suite.sp = serviceprovider.New(cf)
	suite.containerName = "cadvisor"
}

func (suite *PrometheusExprTestSuite) TearDownSuite() {
}

func TestPrometheusExprSuite(t *testing.T) {
	if _, defined := os.LookupEnv("TEST_PROMETHEUS"); !defined {
		t.SkipNow()
		return
	}
	suite.Run(t, new(PrometheusExprTestSuite))
}

func (suite *PrometheusExprTestSuite) TestGetElements() {
	expression := Expression{}
	expression.Metrics = []string{"kube_pod_container_info"}
	expression.QueryLabels = map[string]string{}
	expression.QueryLabels["namespace"] = "vortex"

	str, err := basicExpr(expression)
	results, err := query(suite.sp, str)
	suite.NoError(err)
	suite.NotEqual(0, float32(results[0].Value))

	// Get nil if the result is empty
	results, _ = query(suite.sp, "")
	suite.Equal(model.Vector(nil), results)
}
