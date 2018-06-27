package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PrometheusTestSuite struct {
	suite.Suite
	wc  *restful.Container
	api prometheus.API
}

func (suite *PrometheusTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	//init session
	suite.api = sp.Prometheus.API

	//init restful container
	suite.wc = restful.NewContainer()
	service := newMonitoringService(sp)
	suite.wc.Add(service)
}

func TestPrometheusTestSuite(t *testing.T) {
	suite.Run(t, new(PrometheusTestSuite))
}

func (suite *PrometheusTestSuite) TearDownSuite() {
}

func (suite *PrometheusTestSuite) TestQueryMetrics() {

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/query?query=prometheus_build_info", nil)
	assert.NoError(suite.T(), err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestQueryWrongMetrics() {

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/query?query=!@#$", nil)
	assert.NoError(suite.T(), err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
}
