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

func (suite *UserTestSuite) TestSignUpUser() {
	user := entity.User{
		ID: bson.NewObjectId(),
		LoginCredential: entity.LoginCredential{
			Email:    namesgenerator.GetRandomName(0) + "@linkernetworks.com",
			Password: "p@ssw0rd",
		},
		Username:    "John Doe",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "0900000000",
	}

	bodyBytes, err := json.MarshalIndent(user, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/users/signup", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusCreated, httpWriter)
	defer suite.session.Remove(entity.UserCollectionName, "loginCredential.email", user.LoginCredential.Email)

	// load data to check
	retUser := entity.User{}
	err = suite.session.FindOne(entity.UserCollectionName, bson.M{"loginCredential.email": user.LoginCredential.Email}, &retUser)
	suite.NoError(err)
	suite.NotEqual("", retUser.ID)
	suite.Equal(user.Username, retUser.Username)
	suite.Equal(user.LoginCredential.Email, retUser.LoginCredential.Email)
	// sign up always get the role of user
	suite.Equal("user", retUser.Role)
}

func (suite *UserTestSuite) TestSignUpFailedUser() {
	sameEmail := namesgenerator.GetRandomName(0) + "@linkernetworks.com"
	// given a user already in mongodb
	userFirst := entity.User{
		ID: bson.NewObjectId(),
		LoginCredential: entity.LoginCredential{
			Email:    sameEmail,
			Password: "p@ssw0rd",
		},
		Username:    "John Doe",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "0900000000",
	}
	suite.session.Insert(entity.UserCollectionName, &userFirst)

	userSecond := entity.User{
		ID: bson.NewObjectId(),
		LoginCredential: entity.LoginCredential{
			Email:    sameEmail,
			Password: "p@ssw0rd",
		},
		Username:    "John Doe",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "0900000000",
	}

	bodyBytes, err := json.MarshalIndent(userSecond, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/users/signup", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusConflict, httpWriter)
	defer suite.session.Remove(entity.UserCollectionName, "loginCredential.email", sameEmail)
}

func (suite *UserTestSuite) TestSignInUser() {
	// given a user already in signup
	userCred := entity.LoginCredential{
		Email:    namesgenerator.GetRandomName(0) + "@linkernetworks.com",
		Password: "p@ssw0rd",
	}
	user := entity.User{
		ID:              bson.NewObjectId(),
		LoginCredential: userCred,
		Username:        "John Doe",
		FirstName:       "John",
		LastName:        "Doe",
		PhoneNumber:     "0900000000",
	}

	bodyBytes, err := json.MarshalIndent(user, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/users/signup", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusCreated, httpWriter)
	defer suite.session.Remove(entity.UserCollectionName, "loginCredential.email", user.LoginCredential.Email)

	// do Sign In
	bodyBytesSignIn, err := json.MarshalIndent(userCred, "", "  ")
	suite.NoError(err)

	bodyReaderSignIn := strings.NewReader(string(bodyBytesSignIn))
	httpRequestSignIn, err := http.NewRequest("POST", "http://localhost:7890/v1/users/signin", bodyReaderSignIn)
	suite.NoError(err)

	httpRequestSignIn.Header.Add("Content-Type", "application/json")
	httpWriterSignIn := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriterSignIn, httpRequestSignIn)
	assertResponseCode(suite.T(), http.StatusOK, httpWriterSignIn)
}

func (suite *UserTestSuite) TestSignInFailedUser() {
	userCred := entity.LoginCredential{
		Email:    namesgenerator.GetRandomName(0) + "@linkernetworks.com",
		Password: "p@ssw0rd",
	}

	bodyBytes, err := json.MarshalIndent(userCred, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/users/signin", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusForbidden, httpWriter)
}

func (suite *UserTestSuite) TestCreateUser() {
	user := entity.User{
		ID: bson.NewObjectId(),
		LoginCredential: entity.LoginCredential{
			Email:    namesgenerator.GetRandomName(0) + "@linkernetworks.com",
			Password: "p@ssw0rd",
		},
		Username:    "John Doe",
		Role:        "root",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "0911111111",
		CreatedAt:   &time.Time{},
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
	defer suite.session.Remove(entity.UserCollectionName, "loginCredential.email", user.LoginCredential.Email)

	// load data to check
	retUser := entity.User{}
	err = suite.session.FindOne(entity.UserCollectionName, bson.M{"loginCredential.email": user.LoginCredential.Email}, &retUser)
	suite.NoError(err)
	suite.NotEqual("", retUser.ID)
	suite.Equal(user.Username, retUser.Username)
	suite.Equal(user.LoginCredential.Email, retUser.LoginCredential.Email)

	// We use the new write but empty input which will cause the readEntity Error
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
	// Create again and it should fail since the email exist
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
		ID: bson.NewObjectId(),
		LoginCredential: entity.LoginCredential{
			Email:    "hello@linkernetworks.com",
			Password: "p@ssw0rd",
		},
		Role:        "root",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "",
		CreatedAt:   &time.Time{},
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
		ID: bson.NewObjectId(),
		LoginCredential: entity.LoginCredential{
			Email:    "hello@linkernetworks.com",
			Password: "p@ssw0rd",
		},
		Username:    "John Doe",
		Role:        "root",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "0912121212",
		CreatedAt:   &time.Time{},
	}

	suite.session.Insert(entity.UserCollectionName, &user)
	defer suite.session.Remove(entity.UserCollectionName, "loginCredential.email", user.LoginCredential.Email)

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
	user := entity.User{
		ID: bson.NewObjectId(),
		LoginCredential: entity.LoginCredential{
			Email:    namesgenerator.GetRandomName(0) + "@linkernetworks.com",
			Password: "p@ssw0rd",
		},
		Username:    "John Doe",
		Role:        "root",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "091313l313",
		CreatedAt:   &time.Time{},
	}
	// Create data into mongo manually
	suite.session.Insert(entity.UserCollectionName, &user)
	defer suite.session.Remove(entity.UserCollectionName, "loginCredential.email", user.LoginCredential.Email)

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/users/"+user.ID.Hex(), nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	retUser := entity.User{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &retUser)
	suite.NoError(err)
	suite.Equal(user.Username, retUser.Username)
	suite.Equal(user.LoginCredential.Email, retUser.LoginCredential.Email)
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
	users := []entity.User{}
	count := 3
	for i := 0; i < count; i++ {
		users = append(users, entity.User{
			ID: bson.NewObjectId(),
			LoginCredential: entity.LoginCredential{
				Email:    namesgenerator.GetRandomName(0) + "@linkernetworks.com",
				Password: "p@ssw0rd",
			},
			Username:    "John Doe",
			Role:        "root",
			FirstName:   "John",
			LastName:    "Doe",
			PhoneNumber: "0914141414",
			CreatedAt:   &time.Time{},
		})
	}

	for _, u := range users {
		err := suite.session.Insert(entity.UserCollectionName, &u)
		suite.NoError(err)
		defer suite.session.Remove(entity.UserCollectionName, "loginCredential.email", u.LoginCredential.Email)
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
				suite.Equal(users[i].LoginCredential.Email, u.LoginCredential.Email)
			}
		})
	}
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
