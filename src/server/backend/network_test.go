package backend

import (
	"math/rand"
	"testing"
	"time"

	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"

	"gopkg.in/mgo.v2/bson"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type NetworkTestSuite struct {
	suite.Suite
	session *mongo.Session
	sp      *serviceprovider.Container
}

func (suite *NetworkTestSuite) SetupSuite() {
	cf := config.MustRead("../../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	suite.sp = sp
	// init session
	suite.session = sp.Mongo.NewSession()
}

func (suite *NetworkTestSuite) TearDownSuite() {}

func TestNetworkSuite(t *testing.T) {
	suite.Run(t, new(NetworkTestSuite))
}

func (suite *NetworkTestSuite) TestFindNetworkByID() {
	tName := namesgenerator.GetRandomName(0)
	network := entity.Network{
		OwnerID:    bson.NewObjectId(),
		Type:       entity.FakeNetworkType,
		IsDPDKPort: true, //for fake network, true means success,
		Name:       tName,
		VlanTags:   []int32{},
		BridgeName: namesgenerator.GetRandomName(0),
		Nodes: []entity.Node{
			entity.Node{
				Name:          tName,
				PhyInterfaces: []entity.PhyInterface{},
			},
		},
		CreatedAt: &time.Time{},
	}

	suite.session.Insert(entity.NetworkCollectionName, &network)
	defer suite.session.Remove(entity.NetworkCollectionName, "name", tName)

	retNetwork := entity.Pod{}
	err := suite.session.FindOne(entity.NetworkCollectionName, bson.M{"name": network.Name}, &retNetwork)
	suite.NoError(err)

	network, err = FindNetworkByID(suite.session, retNetwork.ID)
	suite.NoError(err)

	network, err = FindNetworkByID(suite.session, "nonono")
	suite.Error(err)
}
