package prometheuscontroller

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/linkernetworks/vortex/src/serviceprovider"
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

func (suite *PrometheusExprTestSuite) TestBasicExpr() {
	expression := Expression{}
	expression.Metrics = []string{"a", "b", "c"}

	str := basicExpr(expression.Metrics)
	suite.Equal(`__name__=~"a|b|c"`, str)

	str = basicExpr(nil)
	suite.Equal("", str)
}

func (suite *PrometheusExprTestSuite) TestQueryExpr() {
	queryLabels := map[string]string{}
	queryLabels["aKey"] = "aValue"
	str := queryExpr("TEST_STRING", queryLabels)
	suite.Equal(`{TEST_STRING,aKey=~"aValue"}`, str)

	str = queryExpr("TEST_STRING", nil)
	suite.Equal("TEST_STRING", str)
}

func (suite *PrometheusExprTestSuite) TestSumByExpr() {
	sumByLabels := []string{"a", "b", "c"}
	str := sumByExpr("TEST_STRING", sumByLabels)
	suite.Equal(`sum by(a,b,c)(TEST_STRING)`, str)
}

func (suite *PrometheusExprTestSuite) TestSumExpr() {
	str := sumExpr("TEST_STRING")
	suite.Equal(`sum(TEST_STRING)`, str)
}

func (suite *PrometheusExprTestSuite) TestDurationExpr() {
	var duration = "1h"
	str := durationExpr("TEST_STRING", duration)
	suite.Equal(`TEST_STRING[1h]`, str)
}

func (suite *PrometheusExprTestSuite) TestRateExpr() {
	var str = rateExpr("TEST_STRING")
	suite.Equal(`rate(TEST_STRING)`, str)
}

func (suite *PrometheusExprTestSuite) TestEqualExpr() {
	var value = 1234.567
	str := equalExpr("TEST_STRING", value)
	suite.Equal(`TEST_STRING==1234.567`, str)
}

func (suite *PrometheusExprTestSuite) TestMultiplyExpr() {
	value := 1234.567
	str := multiplyExpr("TEST_STRING", value)
	suite.Equal(`TEST_STRING*1234.567`, str)
}
