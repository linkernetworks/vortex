package server

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/stretchr/testify/suite"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type RegistryTestSuite struct {
	suite.Suite
	sp      *serviceprovider.Container
	wc      *restful.Container
	session *mongo.Session
}

func (suite *RegistryTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	suite.sp = sp
	// init session
	suite.session = sp.Mongo.NewSession()
	// init restful container
	suite.wc = restful.NewContainer()
	service := newRegistryService(suite.sp)
	suite.wc.Add(service)
}

func (suite *RegistryTestSuite) TearDownSuite() {}

func TestRegistrySuite(t *testing.T) {
	suite.Run(t, new(RegistryTestSuite))
}

func (suite *RegistryTestSuite) TestRegistryBasicAuth() {
	cred := entity.RegistryBasicAuthCredential{
		Username: "test",
		Password: "test",
	}

	bodyBytes, err := json.MarshalIndent(cred, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/registry/auth", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusUnauthorized, httpWriter)
}
