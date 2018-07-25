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
	wc      *restful.Container
	session *mongo.Session
}

func (suite *UserTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	suite.sp = sp
	//init session
	suite.session = sp.Mongo.NewSession()
	//init restful container
	suite.wc = restful.NewContainer()
	user := newUserService(suite.sp)
	suite.wc.Add(user)
}

func (suite *UserTestSuite) TearDownSuite() {}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (suite *UserTestSuite) TestCreateUser() {
	mailAccount := namesgenerator.GetRandomName(0)
	user := entity.User{
		ID:          bson.NewObjectId(),
		Username:    "John Doe",
		Email:       mailAccount + "@linkernetworks.com",
		Password:    "p@ssw0rd",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "0911123456",
	}

	bodyBytes, err := json.MarshalIndent(user, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/users", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
	defer suite.session.Remove(entity.UserCollectionName, "email", user.Email)

	// load data to check
	retUser := entity.User{}
	err = suite.session.FindOne(entity.UserCollectionName, bson.M{"email": user.Email}, &retUser)
	suite.NoError(err)
	suite.NotEqual("", retUser.ID)
	suite.Equal(user.Username, retUser.Username)
	suite.Equal(user.Email, retUser.Email)

	// We use the new write but empty input which will cause the readEntity Error
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
	// Create again and it should fail since the name exist
	bodyReader = strings.NewReader(string(bodyBytes))
	httpRequest, err = http.NewRequest("POST", "http://localhost:7890/v1/users", bodyReader)
	suite.NoError(err)
	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusConflict, httpWriter)
}

func (suite *UserTestSuite) TestCreateUserFail() {
	user := entity.User{
		ID:          bson.NewObjectId(),
		Username:    "John Doe",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "0911123456",
	}

	bodyBytes, err := json.MarshalIndent(user, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/users", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
}

func (suite *UserTestSuite) TestDeleteUser() {
	user := entity.User{
		ID:          bson.NewObjectId(),
		Username:    "John Doe",
		Email:       "hello@linkernetworks.com",
		Password:    "p@ssw0rd",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "0911123456",
	}

	err := suite.session.Insert(entity.UserCollectionName, &user)
	suite.NoError(err)

	bodyBytes, err := json.MarshalIndent(user, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/users/"+user.ID.Hex(), bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	n, err := suite.session.Count(entity.UserCollectionName, bson.M{"_id": user.ID})
	suite.NoError(err)
	suite.Equal(0, n)
}

func (suite *UserTestSuite) TestDeleteUserWithInvalidID() {
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/users/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
}

// For Get/List, we only return mongo document
func (suite *UserTestSuite) TestGetUser() {
	mailAccount := namesgenerator.GetRandomName(0)
	user := entity.User{
		ID:          bson.NewObjectId(),
		Username:    "John Doe",
		Email:       mailAccount + "@linkernetworks.com",
		Password:    "p@ssw0rd",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "0911123456",
	}
	// Create data into mongo manually
	suite.session.C(entity.UserCollectionName).Insert(user)
	defer suite.session.Remove(entity.UserCollectionName, "email", user.Email)

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/users/"+user.ID.Hex(), nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	retUser := entity.User{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &retUser)
	suite.NoError(err)
	suite.Equal(user.Username, retUser.Username)
	suite.Equal(user.Email, retUser.Email)
}

func (suite *UserTestSuite) TestGetUserWithInvalidID() {
	// Get data with non-exits ID
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/users/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusNotFound, httpWriter)
}

func (suite *UserTestSuite) TestListUser() {
	/*
		users := []entity.User{}
		count := 3
		for i := 0; i < count; i++ {
			mailAccount := namesgenerator.GetRandomName(i)
			users = append(users, entity.User{
				ID:          bson.NewObjectId(),
				Username:    "John Doe",
				Email:       mailAccount + "@linkernetworks.com",
				Password:    "p@ssw0rd",
				FirstName:   "John",
				LastName:    "Doe",
				PhoneNumber: "0911123456",
			})
		}

		for _, u := range users {
			suite.session.C(entity.UserCollectionName).Insert(u)
			defer suite.session.Remove(entity.UserCollectionName, "email", u.Email)
		}

		testCases := []struct {
			page       string
			pageSize   string
			expectSize int
		}{
			{"", "", count},
			{"1", "1", count},
			{"1", "3", count},
		}

		for _, tc := range testCases {
			caseName := "page:pageSize" + tc.page + ":" + tc.pageSize
			suite.T().Run(caseName, func(t *testing.T) {
				// list data by default page and page_size
				url := "http://localhost:7890/v1/users/"
				if tc.page != "" || tc.pageSize != "" {
					url = "http://localhost:7890/v1/users?"
					url += "page=" + tc.page + "%" + "page_size" + tc.pageSize
				}
				httpRequest, err := http.NewRequest("GET", url, nil)
				suite.NoError(err)

				httpWriter := httptest.NewRecorder()
				suite.wc.Dispatch(httpWriter, httpRequest)
				assertResponseCode(suite.T(), http.StatusOK, httpWriter)

				retUsers := []entity.User{}
				err = json.Unmarshal(httpWriter.Body.Bytes(), &retUsers)
				suite.NoError(err)
				suite.Equal(tc.expectSize, len(retUsers))
				for i, u := range retUsers {
					suite.Equal(users[i].Username, u.Username)
					suite.Equal(users[i].Email, u.Email)
				}
			})
		}*/
}

func (suite *UserTestSuite) TestListUserWithInvalidPage() {
	// Get data with non-exits ID
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/users?page=asdd", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/users?page_size=asdd", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/users?page=-1", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusInternalServerError, httpWriter)
}
