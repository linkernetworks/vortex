package namespace

import (
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type NamespaceTestSuite struct {
	suite.Suite
	sp *serviceprovider.Container
}

func (suite *NamespaceTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	suite.sp = serviceprovider.NewForTesting(cf)
}

func (suite *NamespaceTestSuite) TearDownSuite() {
}

func TestNamespaceSuite(t *testing.T) {
	suite.Run(t, new(NamespaceTestSuite))
}

func (suite *NamespaceTestSuite) TestCreateDeleteNamespace() {
	namespaceName := namesgenerator.GetRandomName(0)
	namespace := &entity.Namespace{
		ID:   bson.NewObjectId(),
		Name: namespaceName,
	}

	err := CreateNamespace(suite.sp, namespace)
	suite.NoError(err)

	err = DeleteNamespace(suite.sp, namespace)
	suite.NoError(err)
}
