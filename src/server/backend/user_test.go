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

type UserTestSuite struct {
	suite.Suite
	sp      *serviceprovider.Container
	session *mongo.Session
}

func (suite *UserTestSuite) SetupSuite() {
	cf := config.MustRead("../../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	suite.sp = sp
	// init session
	suite.session = sp.Mongo.NewSession()
}

func (suite *UserTestSuite) TearDownSuite() {}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (suite *UserTestSuite) TestFindUserByID() {
	user := entity.User{
		ID: bson.NewObjectId(),
		LoginCredential: entity.LoginCredential{
			Username: namesgenerator.GetRandomName(0) + "@linkernetworks.com",
			Password: "p@ssw0rd",
		},
		DisplayName: "John Doe",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "0900000000",
	}

	suite.session.Insert(entity.UserCollectionName, &user)
	defer suite.session.Remove(entity.UserCollectionName, "loginCredential.username", user.LoginCredential.Username)

	// load data to check
	retUser := entity.User{}
	err := suite.session.FindOne(entity.UserCollectionName, bson.M{"loginCredential.username": user.LoginCredential.Username}, &retUser)
	suite.NoError(err)
	suite.NotEqual("", retUser.ID)
	suite.Equal(user.DisplayName, retUser.DisplayName)
	suite.Equal(user.LoginCredential.Username, retUser.LoginCredential.Username)

	user, err = FindUserByID(suite.session, retUser.ID)
	suite.NoError(err)

	user, err = FindUserByID(suite.session, "nonono")
	suite.Error(err)
}
