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
	ns "github.com/linkernetworks/vortex/src/namespace"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type NamespaceTestSuite struct {
	suite.Suite
	sp      *serviceprovider.Container
	wc      *restful.Container
	session *mongo.Session
}

func (suite *NamespaceTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	suite.sp = sp
	//init session
	suite.session = sp.Mongo.NewSession()
	//init restful container
	suite.wc = restful.NewContainer()
	service := newNamespaceService(suite.sp)
	suite.wc.Add(service)
}

func (suite *NamespaceTestSuite) TearDownSuite() {
}

func TestNamespaceSuite(t *testing.T) {
	suite.Run(t, new(NamespaceTestSuite))
}

func (suite *NamespaceTestSuite) TestCreateNamespace() {

	nsName := namesgenerator.GetRandomName(0)
	namespace := entity.Namespace{
		ID:     bson.NewObjectId(),
		Name:   nsName,
		Labels: map[string]string{},
	}

	bodyBytes, err := json.MarshalIndent(namespace, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/namespaces", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusCreated, httpWriter)
	defer suite.session.Remove(entity.NamespaceCollectionName, "name", namespace.Name)

	//load data to check
	retNamespace := entity.Namespace{}
	err = suite.session.FindOne(entity.NamespaceCollectionName, bson.M{"name": namespace.Name}, &retNamespace)
	suite.NoError(err)
	suite.NotEqual("", retNamespace.ID)
	suite.Equal(namespace.Name, retNamespace.Name)

	//We use the new write but empty input which will cause the readEntity Error
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	//Create again and it should fail since the name exist
	bodyReader = strings.NewReader(string(bodyBytes))
	httpRequest, err = http.NewRequest("POST", "http://localhost:7890/v1/namespaces", bodyReader)
	suite.NoError(err)
	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusConflict, httpWriter)

	err = ns.DeleteNamespace(suite.sp, &namespace)
	suite.NoError(err)
}

func (suite *NamespaceTestSuite) TestCreateNamespaceFail() {
	nsName := namesgenerator.GetRandomName(0)
	namespace := entity.Namespace{
		ID:   bson.NewObjectId(),
		Name: nsName,
	}

	bodyBytes, err := json.MarshalIndent(namespace, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/namespaces", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
}

func (suite *NamespaceTestSuite) TestDeleteNamespace() {
	nsName := namesgenerator.GetRandomName(0)
	namespace := entity.Namespace{
		ID:   bson.NewObjectId(),
		Name: nsName,
	}

	err := ns.CreateNamespace(suite.sp, &namespace)
	suite.NoError(err)

	err = suite.session.Insert(entity.NamespaceCollectionName, &namespace)
	suite.NoError(err)

	bodyBytes, err := json.MarshalIndent(namespace, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/namespaces/"+namespace.ID.Hex(), bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	n, err := suite.session.Count(entity.NamespaceCollectionName, bson.M{"_id": namespace.ID})
	suite.NoError(err)
	suite.Equal(0, n)
}

func (suite *NamespaceTestSuite) TestDeleteNamespaceWithInvalidID() {
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/namespaces/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
}

//For Get/List, we only return mongo document
func (suite *NamespaceTestSuite) TestGetNamespace() {
	nsName := namesgenerator.GetRandomName(0)
	namespace := entity.Namespace{
		ID:   bson.NewObjectId(),
		Name: nsName,
	}

	//Create data into mongo manually
	suite.session.C(entity.NamespaceCollectionName).Insert(namespace)
	defer suite.session.Remove(entity.NamespaceCollectionName, "name", nsName)

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/namespaces/"+namespace.ID.Hex(), nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	namespace = entity.Namespace{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &namespace)
	suite.NoError(err)
	suite.Equal(nsName, namespace.Name)
}

func (suite *NamespaceTestSuite) TestGetNamespaceWithInvalidID() {
	//Get data with non-exits ID
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/namespaces/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusNotFound, httpWriter)
}

func (suite *NamespaceTestSuite) TestListNamespace() {
	namespaces := []entity.Namespace{}
	count := 3
	for i := 0; i < count; i++ {
		namespaces = append(namespaces, entity.Namespace{
			ID:   bson.NewObjectId(),
			Name: namesgenerator.GetRandomName(0),
		})
	}

	for _, n := range namespaces {
		suite.session.C(entity.NamespaceCollectionName).Insert(n)
		defer suite.session.Remove(entity.NamespaceCollectionName, "_id", n.ID)
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
			//list data by default page and page_size
			url := "http://localhost:7890/v1/namespaces/"
			if tc.page != "" || tc.pageSize != "" {
				url = "http://localhost:7890/v1/namespaces?"
				url += "page=" + tc.page + "%" + "page_size" + tc.pageSize
			}
			httpRequest, err := http.NewRequest("GET", url, nil)
			suite.NoError(err)

			httpWriter := httptest.NewRecorder()
			suite.wc.Dispatch(httpWriter, httpRequest)
			assertResponseCode(suite.T(), http.StatusOK, httpWriter)

			retNamespaces := []entity.Namespace{}
			err = json.Unmarshal(httpWriter.Body.Bytes(), &retNamespaces)
			suite.NoError(err)
			suite.Equal(tc.expectSize, len(retNamespaces))
			for i, n := range retNamespaces {
				suite.Equal(namespaces[i].Name, n.Name)
			}
		})
	}
}

func (suite *NamespaceTestSuite) TestListNamespaceWithInvalidPage() {
	//Get data with non-exits ID
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/namespaces?page=asdd", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/namespaces?page_size=asdd", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/namespaces?page=-1", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusInternalServerError, httpWriter)
}
