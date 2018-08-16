package server

import (
	_ "encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"

	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	_ "github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"
	_ "gopkg.in/mgo.v2/bson"
	//corev1 "k8s.io/api/core/v1"

	"testing"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type OVSTestSuite struct {
	suite.Suite
	sp      *serviceprovider.Container
	wc      *restful.Container
	session *mongo.Session
	storage entity.Storage
}

func (suite *OVSTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	suite.sp = sp
	//init session
	suite.session = sp.Mongo.NewSession()
	//init restful container
	suite.wc = restful.NewContainer()
	service := newOVSService(suite.sp)
	suite.wc.Add(service)
}

func (suite *OVSTestSuite) TearDownSuite() {
}

func TestOVSSuite(t *testing.T) {
	suite.Run(t, new(OVSTestSuite))
}

func (suite *OVSTestSuite) TestGetOVSPortStatsFail() {
	//Empty data
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/ovs/portstat", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/ovs/portstat?nodeName=11", nil)
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/ovs/portstat?nodeName=11&&bridgeName=111", nil)
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusInternalServerError, httpWriter)
}
