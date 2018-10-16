package server

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type AppSuite struct {
	suite.Suite
	app App
}

func (suite *AppSuite) SetupSuite() {
}

func (suite *AppSuite) TearDownSuite() {
}

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppSuite))
}

func (suite *AppSuite) TestLoadConfig() {
	a := suite.app.LoadConfig("../../config/testing.json")
	suite.NotNil(a)
}

func (suite *AppSuite) InitilizeService() {
	suite.app.InitilizeService()
	suite.NotNil(suite.app.ServiceProvider)
}
