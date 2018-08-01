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
	p "github.com/linkernetworks/vortex/src/pod"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type PodTestSuite struct {
	suite.Suite
	sp      *serviceprovider.Container
	wc      *restful.Container
	session *mongo.Session
}

func (suite *PodTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	suite.sp = sp
	//init session
	suite.session = sp.Mongo.NewSession()
	//init restful container
	suite.wc = restful.NewContainer()
	service := newPodService(suite.sp)
	suite.wc.Add(service)
}

func (suite *PodTestSuite) TearDownSuite() {}

func TestPodSuite(t *testing.T) {
	suite.Run(t, new(PodTestSuite))
}

func (suite *PodTestSuite) TestCreatePod() {
	namespace := "default"
	containers := []entity.Container{
		{
			Name:    namesgenerator.GetRandomName(0),
			Image:   "busybox",
			Command: []string{"sleep", "3600"},
		},
	}
	tName := namesgenerator.GetRandomName(0)
	pod := entity.Pod{
		Name:          tName,
		Namespace:     namespace,
		Labels:        map[string]string{},
		Containers:    containers,
		Volumes:       []entity.PodVolume{},
		Networks:      []entity.PodNetwork{},
		Capability:    true,
		RestartPolicy: "Never",
		NetworkType:   entity.PodHostNetwork,
		NodeAffinity:  []string{},
	}
	bodyBytes, err := json.MarshalIndent(pod, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/pods", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusCreated, httpWriter)
	defer suite.session.Remove(entity.PodCollectionName, "name", pod.Name)

	//load data to check
	retPod := entity.Pod{}
	err = suite.session.FindOne(entity.PodCollectionName, bson.M{"name": pod.Name}, &retPod)
	suite.NoError(err)
	suite.NotEqual("", retPod.ID)
	suite.Equal(pod.Name, retPod.Name)
	suite.Equal(len(pod.Containers), len(retPod.Containers))

	//We use the new write but empty input which will cause the readEntity Error
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
	//Create again and it should fail since the name exist
	bodyReader = strings.NewReader(string(bodyBytes))
	httpRequest, err = http.NewRequest("POST", "http://localhost:7890/v1/pods", bodyReader)
	suite.NoError(err)
	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusInternalServerError, httpWriter)

	err = p.DeletePod(suite.sp, &retPod)
	suite.NoError(err)
}

func (suite *PodTestSuite) TestCreatePodFail() {
	namespace := "default"
	containers := []entity.Container{
		{
			Name:    namesgenerator.GetRandomName(0),
			Image:   "busybox",
			Command: []string{"sleep", "3600"},
		},
	}
	tName := namesgenerator.GetRandomName(0)
	pod := entity.Pod{
		Name:       tName,
		Namespace:  namespace,
		Containers: containers,
		Volumes: []entity.PodVolume{
			{Name: namesgenerator.GetRandomName(0)},
		},
	}

	bodyBytes, err := json.MarshalIndent(pod, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/pods", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
}

func (suite *PodTestSuite) TestDeletePod() {
	namespace := "default"
	containers := []entity.Container{
		{
			Name:    namesgenerator.GetRandomName(0),
			Image:   "busybox",
			Command: []string{"sleep", "3600"},
		},
	}
	tName := namesgenerator.GetRandomName(0)
	pod := entity.Pod{
		ID:            bson.NewObjectId(),
		Name:          tName,
		Namespace:     namespace,
		Containers:    containers,
		Capability:    true,
		RestartPolicy: "Never",
		NetworkType:   entity.PodHostNetwork,
		NodeAffinity:  []string{},
	}

	err := p.CreatePod(suite.sp, &pod)
	suite.NoError(err)

	err = suite.session.Insert(entity.PodCollectionName, &pod)
	suite.NoError(err)

	bodyBytes, err := json.MarshalIndent(pod, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/pods/"+pod.ID.Hex(), bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	n, err := suite.session.Count(entity.PodCollectionName, bson.M{"_id": pod.ID})
	suite.NoError(err)
	suite.Equal(0, n)
}

func (suite *PodTestSuite) TestDeletePodWithInvalidID() {
	httpRequest, err := http.NewRequest("DELETE", "http://localhost:7890/v1/pods/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)
}

//For Get/List, we only return mongo document
func (suite *PodTestSuite) TestGetPod() {
	namespace := "default"
	containers := []entity.Container{
		{
			Name:    namesgenerator.GetRandomName(0),
			Image:   "busybox",
			Command: []string{"sleep", "3600"},
		},
	}
	tName := namesgenerator.GetRandomName(0)
	pod := entity.Pod{
		ID:         bson.NewObjectId(),
		Name:       tName,
		Namespace:  namespace,
		Containers: containers,
	}

	//Create data into mongo manually
	suite.session.C(entity.PodCollectionName).Insert(pod)
	defer suite.session.Remove(entity.PodCollectionName, "name", tName)

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/pods/"+pod.ID.Hex(), nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	pod = entity.Pod{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &pod)
	suite.NoError(err)
	suite.Equal(tName, pod.Name)
	suite.Equal(len(containers), len(pod.Containers))
}

func (suite *PodTestSuite) TestGetPodWithInvalidID() {
	//Get data with non-exits ID
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/pods/"+bson.NewObjectId().Hex(), nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusNotFound, httpWriter)
}

func (suite *PodTestSuite) TestListPod() {
	namespace := "default"
	pods := []entity.Pod{}
	count := 3
	for i := 0; i < count; i++ {
		containers := []entity.Container{
			{
				Name:    namesgenerator.GetRandomName(0),
				Image:   "busybox",
				Command: []string{"sleep", "3600"},
			},
		}
		pods = append(pods, entity.Pod{
			ID:         bson.NewObjectId(),
			Name:       namesgenerator.GetRandomName(0),
			Namespace:  namespace,
			Containers: containers,
		})
	}

	for _, p := range pods {
		suite.session.C(entity.PodCollectionName).Insert(p)
		defer suite.session.Remove(entity.PodCollectionName, "_id", p.ID)
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
			url := "http://localhost:7890/v1/pods/"
			if tc.page != "" || tc.pageSize != "" {
				url = "http://localhost:7890/v1/pods?"
				url += "page=" + tc.page + "%" + "page_size" + tc.pageSize
			}
			httpRequest, err := http.NewRequest("GET", url, nil)
			suite.NoError(err)

			httpWriter := httptest.NewRecorder()
			suite.wc.Dispatch(httpWriter, httpRequest)
			assertResponseCode(suite.T(), http.StatusOK, httpWriter)

			retPods := []entity.Pod{}
			err = json.Unmarshal(httpWriter.Body.Bytes(), &retPods)
			suite.NoError(err)
			suite.Equal(tc.expectSize, len(retPods))
			for i, p := range retPods {
				suite.Equal(pods[i].Name, p.Name)
				suite.Equal(len(pods[i].Containers), len(p.Containers))
			}
		})
	}
}

func (suite *PodTestSuite) TestListPodWithInvalidPage() {
	//Get data with non-exits ID
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/pods?page=asdd", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/pods?page_size=asdd", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusBadRequest, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/pods?page=-1", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusInternalServerError, httpWriter)
}
