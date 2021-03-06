package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	svc "github.com/linkernetworks/vortex/src/service"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type ServiceTestSuite struct {
	suite.Suite
	sp        *serviceprovider.Container
	wc        *restful.Container
	session   *mongo.Session
	JWTBearer string
}

func (suite *ServiceTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	suite.sp = sp
	// init session
	suite.session = sp.Mongo.NewSession()
	// init restful container
	suite.wc = restful.NewContainer()

	serviceService := newServiceService(suite.sp)
	userService := newUserService(suite.sp)

	suite.wc.Add(serviceService)
	suite.wc.Add(userService)

	token, _ := loginGetToken(suite.wc)
	suite.NotEmpty(token)
	suite.JWTBearer = "Bearer " + token
}

func (suite *ServiceTestSuite) TearDownSuite() {}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) TestCreateService() {
	selector := map[string]string{
		"podname": "awesome",
	}
	ports := []entity.ServicePort{
		{
			Name:       namesgenerator.GetRandomName(0),
			Port:       int32(80),
			TargetPort: 80,
			NodePort:   int32(30000),
		},
	}

	serviceName := namesgenerator.GetRandomName(0)
	service := entity.Service{
		OwnerID:   bson.NewObjectId(),
		ID:        bson.NewObjectId(),
		Name:      serviceName,
		Namespace: "default",
		Type:      "NodePort",
		Selector:  selector,
		Ports:     ports,
	}

	bodyBytes, err := json.MarshalIndent(service, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/services", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusCreated, httpWriter)
	defer suite.session.Remove(entity.ServiceCollectionName, "name", service.Name)

	//load data to check
	retService := entity.Service{}
	err = suite.session.FindOne(entity.ServiceCollectionName, bson.M{"name": service.Name}, &retService)
	suite.NoError(err)
	suite.NotEqual("", retService.ID)
	suite.Equal(service.Name, retService.Name)
	suite.Equal(len(service.Ports), len(retService.Ports))

	//We use the new write but empty input which will cause the readEntity Error
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	//Create again and it should fail since the name exist
	bodyReader = strings.NewReader(string(bodyBytes))
	httpRequest, err = http.NewRequest("POST", "http://localhost:7890/v1/services", bodyReader)
	suite.NoError(err)
	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusConflict, httpWriter)

	err = svc.DeleteService(suite.sp, &service)
	suite.NoError(err)
}

func (suite *ServiceTestSuite) TestCreateServiceFail() {
	serviceName := namesgenerator.GetRandomName(0)
	service := entity.Service{
		ID:      bson.NewObjectId(),
		OwnerID: bson.NewObjectId(),
		Name:    serviceName,
	}

	bodyBytes, err := json.MarshalIndent(service, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/services", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
}

func (suite *ServiceTestSuite) TestDeleteService() {
	selector := map[string]string{
		"podname": "awesome",
	}
	ports := []entity.ServicePort{
		{
			Name:       namesgenerator.GetRandomName(0),
			Port:       int32(80),
			TargetPort: 80,
			NodePort:   int32(30000),
		},
	}

	serviceName := namesgenerator.GetRandomName(0)
	service := entity.Service{
		ID:        bson.NewObjectId(),
		OwnerID:   bson.NewObjectId(),
		Name:      serviceName,
		Namespace: "default",
		Type:      "NodePort",
		Selector:  selector,
		Ports:     ports,
	}

	err := svc.CreateService(suite.sp, &service)
	suite.NoError(err)

	err = suite.session.Insert(entity.ServiceCollectionName, &service)
	suite.NoError(err)

	bodyBytes, err := json.MarshalIndent(service, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/services/"+service.ID.Hex(), bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	n, err := suite.session.Count(entity.ServiceCollectionName, bson.M{"_id": service.ID})
	suite.NoError(err)
	suite.Equal(0, n)
}

func (suite *ServiceTestSuite) TestDeleteServiceWithInvalidID() {
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/services/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
}

//For Get/List, we only return mongo document
func (suite *ServiceTestSuite) TestGetService() {
	selector := map[string]string{
		"podname": "awesome",
	}
	ports := []entity.ServicePort{
		{
			Name:       namesgenerator.GetRandomName(0),
			Port:       int32(80),
			TargetPort: 80,
			NodePort:   int32(30000),
		},
	}

	serviceName := namesgenerator.GetRandomName(0)
	service := entity.Service{
		ID:        bson.NewObjectId(),
		OwnerID:   bson.NewObjectId(),
		Name:      serviceName,
		Namespace: "default",
		Type:      "NodePort",
		Selector:  selector,
		Ports:     ports,
	}

	//Create data into mongo manually
	suite.session.C(entity.ServiceCollectionName).Insert(service)
	defer suite.session.Remove(entity.ServiceCollectionName, "name", serviceName)

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/services/"+service.ID.Hex(), nil)
	suite.NoError(err)

	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	service = entity.Service{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &service)
	suite.NoError(err)
	suite.Equal(serviceName, service.Name)
	suite.Equal(len(ports), len(service.Ports))
}

func (suite *ServiceTestSuite) TestGetServiceWithInvalidID() {
	//Get data with non-exits ID
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/services/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)

	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusNotFound, httpWriter)
}

func (suite *ServiceTestSuite) TestListService() {
	services := []entity.Service{}
	count := 3
	for i := 0; i < count; i++ {
		selector := map[string]string{
			"podname": "awesome",
		}
		ports := []entity.ServicePort{
			{
				Name:       namesgenerator.GetRandomName(0),
				Port:       int32(80),
				TargetPort: 80,
				NodePort:   int32(30000),
			},
		}

		services = append(services, entity.Service{
			ID:        bson.NewObjectId(),
			OwnerID:   bson.NewObjectId(),
			Name:      namesgenerator.GetRandomName(0),
			Namespace: "default",
			Type:      "NodePort",
			Selector:  selector,
			Ports:     ports,
		})
	}

	for _, s := range services {
		suite.session.C(entity.ServiceCollectionName).Insert(s)
		defer suite.session.Remove(entity.ServiceCollectionName, "_id", s.ID)
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
			url := "http://localhost:7890/v1/services/"
			if tc.page != "" || tc.pageSize != "" {
				url = "http://localhost:7890/v1/services?"
				url += "page=" + tc.page + "%" + "page_size" + tc.pageSize
			}
			httpRequest, err := http.NewRequest("GET", url, nil)
			suite.NoError(err)

			httpRequest.Header.Add("Authorization", suite.JWTBearer)
			httpWriter := httptest.NewRecorder()
			suite.wc.Dispatch(httpWriter, httpRequest)
			assertResponseCode(suite.T(), http.StatusOK, httpWriter)

			retServices := []entity.Service{}
			err = json.Unmarshal(httpWriter.Body.Bytes(), &retServices)
			suite.NoError(err)
			suite.Equal(tc.expectSize, len(retServices))
			for i, s := range retServices {
				suite.Equal(services[i].Name, s.Name)
				suite.Equal(len(services[i].Ports), len(s.Ports))
			}
		})
	}
}

func (suite *ServiceTestSuite) TestListServiceWithInvalidPage() {
	//Get data with non-exits ID
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/services?page=asdd", nil)
	suite.NoError(err)

	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/services?page_size=asdd", nil)
	suite.NoError(err)

	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/services?page=-1", nil)
	suite.NoError(err)

	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusInternalServerError, httpWriter)
}

func (suite *ServiceTestSuite) TestUploadServiceYAML() {
	filename := "../../testYAMLs/service.yaml"

	bodyBuf := bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(bodyBuf)

	// use the bodyWriter to write the Part headers to the buffer
	_, err := bodyWriter.CreateFormFile("file", filename)
	suite.NoError(err)

	// the file data will be the second part of the body
	file, err := os.Open(filename)
	suite.NoError(err)

	// need to know the boundary to properly close the part myself.
	boundary := bodyWriter.Boundary()
	//close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until
	// writing to the socket buffer.
	requestReader := io.MultiReader(bodyBuf, file, closeBuf)
	fileStat, err := file.Stat()
	suite.NoError(err)

	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/services/upload/yaml", requestReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpRequest.ContentLength = fileStat.Size() + int64(bodyBuf.Len()) + int64(closeBuf.Len())
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	defer suite.session.Remove(entity.ServiceCollectionName, "name", "upload-service")

	assertResponseCode(suite.T(), http.StatusCreated, httpWriter)

	//load data to check
	retService := entity.Service{}
	err = suite.session.FindOne(entity.ServiceCollectionName, bson.M{"name": "upload-service"}, &retService)
	suite.NoError(err)
	suite.NotEqual("", retService.ID)
	suite.Equal("upload-service", retService.Name)

	//Create again and it should fail since the name exist
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusConflict, httpWriter)

	err = svc.DeleteService(suite.sp, &retService)
	suite.NoError(err)
}

func (suite *ServiceTestSuite) TestUploadServiceYAMLFail() {
	filename := "../../testYAMLs/namespace.yaml"

	bodyBuf := bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(bodyBuf)

	// use the bodyWriter to write the Part headers to the buffer
	_, err := bodyWriter.CreateFormFile("file", filename)
	suite.NoError(err)

	// the file data will be the second part of the body
	file, err := os.Open(filename)
	suite.NoError(err)

	// need to know the boundary to properly close the part myself.
	boundary := bodyWriter.Boundary()
	//close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until
	// writing to the socket buffer.
	requestReader := io.MultiReader(bodyBuf, file, closeBuf)
	fileStat, err := file.Stat()
	suite.NoError(err)

	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/services/upload/yaml", requestReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpRequest.ContentLength = fileStat.Size() + int64(bodyBuf.Len()) + int64(closeBuf.Len())
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)

	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
}
