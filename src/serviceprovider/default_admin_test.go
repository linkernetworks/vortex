package serviceprovider

import (
	"testing"

	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/stretchr/testify/suite"
)

type CreateDefaultUserSuite struct {
	suite.Suite
	cf      config.Config
	session *mongo.Session
	service *mongo.Service
}

func (suite *CreateDefaultUserSuite) SetupSuite() {
	suite.cf = config.MustRead("../../config/testing.json")
	sp := NewForTesting(suite.cf)
	suite.service = sp.Mongo
	suite.session = suite.service.NewSession()
}

func (suite *CreateDefaultUserSuite) TearDownSuite() {
	suite.session.Remove(entity.UserCollectionName, "loginCredential.username", "admin@vortex.com")
	defer suite.session.Close()
}

func TestCreateDefaultUserSuite(t *testing.T) {
	suite.Run(t, new(CreateDefaultUserSuite))
}

func (suite *CreateDefaultUserSuite) TestDefaultUserCreate() {
	err := createDefaultUser(suite.service)
	suite.NoError(err)
}
