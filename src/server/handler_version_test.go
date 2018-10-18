package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	restful "github.com/emicklei/go-restful"

	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/stretchr/testify/suite"
)

type VersionTestSuite struct {
	suite.Suite
	wc *restful.Container
	sp *serviceprovider.Container
}

func (suite *VersionTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	//init restful container
	suite.wc = restful.NewContainer()
	service := newVersionService(sp)
	suite.wc.Add(service)
}

func (suite *VersionTestSuite) TearDownSuite() {}

func TestVersionSuite(t *testing.T) {
	suite.Run(t, new(VersionTestSuite))
}

func (suite *VersionTestSuite) TestVersion() {
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/version", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}
