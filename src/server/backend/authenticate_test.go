package backend

import (
	"math/rand"
	"testing"
	"time"

	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type AuthenticateTestSuite struct {
	suite.Suite
	session           *mongo.Session
	sp                *serviceprovider.Container
	plainTextPassword string
}

func (suite *AuthenticateTestSuite) SetupSuite() {
	cf := config.MustRead("../../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	// init session
	suite.session = sp.Mongo.NewSession()
	suite.sp = sp

	suite.plainTextPassword = "@uthentic@te"

	hashedPassword, err := HashPassword(suite.plainTextPassword)
	suite.NoError(err)

	user := entity.User{
		ID: bson.NewObjectId(),
		LoginCredential: entity.LoginCredential{
			Username: "auth@linkernetworks.com",
			Password: hashedPassword,
		},
		DisplayName: "John Doe",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "091111l111",
	}
	err = suite.session.Insert(entity.UserCollectionName, &user)
	suite.NoError(err)
}

func (suite *AuthenticateTestSuite) TearDownSuite() {
	suite.session.Remove(
		entity.UserCollectionName,
		"loginCredential.email",
		"auth@linkernetworks.com",
	)
}

func TestAuthenticateSuite(t *testing.T) {
	suite.Run(t, new(AuthenticateTestSuite))
}

func (suite *AuthenticateTestSuite) TestAuthenticate() {
	CorrectCred := entity.LoginCredential{
		Username: "auth@linkernetworks.com",
		Password: suite.plainTextPassword,
	}
	user, passed, err := Authenticate(suite.session, CorrectCred)
	suite.NoError(err)
	suite.True(passed)
	suite.Equal(CorrectCred.Username, user.LoginCredential.Username)
}

func (suite *AuthenticateTestSuite) TestFailedAuthenticate() {
	WrongCred := entity.LoginCredential{
		Username: "auth@linkernetworks.com",
		Password: "wrongPasswordOX",
	}
	_, passed, err := Authenticate(suite.session, WrongCred)
	suite.NoError(err)
	suite.False(passed)
}
