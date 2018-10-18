package serviceprovider

import (
	"testing"

	"github.com/linkernetworks/vortex/src/config"
	"github.com/stretchr/testify/suite"
)

type ServiceProviderSuite struct {
	suite.Suite
}

func (suite *ServiceProviderSuite) SetupSuite() {
}

func (suite *ServiceProviderSuite) TearDownSuite() {
}

func TestServiceProviderSuite(t *testing.T) {
	suite.Run(t, new(ServiceProviderSuite))
}

func (suite *ServiceProviderSuite) TestContainer() {
	container := NewContainer("../../config/testing.json")
	suite.NotNil(container)
}

func (suite *ServiceProviderSuite) TestNewForTesting() {
	cf := config.MustRead("../../config/testing.json")
	sp := NewForTesting(cf)
	suite.NotNil(sp)
}
