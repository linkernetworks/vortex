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

type PodTestSuite struct {
	suite.Suite
	sp      *serviceprovider.Container
	session *mongo.Session
}

func (suite *PodTestSuite) SetupSuite() {
	cf := config.MustRead("../../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	suite.sp = sp
	// init session
	suite.session = sp.Mongo.NewSession()
}

func (suite *PodTestSuite) TearDownSuite() {}

func TestPodSuite(t *testing.T) {
	suite.Run(t, new(PodTestSuite))
}

func (suite *PodTestSuite) TestFindPodByName() {
	namespace := "default"
	containers := []entity.Container{
		{
			Name:                  namesgenerator.GetRandomName(0),
			Image:                 "busybox",
			Command:               []string{"sleep", "3600"},
			ResourceRequestCPU:    0,
			ResourceRequestMemory: 0,
		},
	}
	tName := namesgenerator.GetRandomName(0)
	pod := entity.Pod{
		OwnerID:       bson.NewObjectId(),
		Name:          tName,
		Namespace:     namespace,
		Labels:        map[string]string{},
		EnvVars:       map[string]string{},
		Containers:    containers,
		Volumes:       []entity.PodVolume{},
		Networks:      []entity.PodNetwork{},
		Capability:    true,
		RestartPolicy: "Never",
		NetworkType:   entity.PodHostNetwork,
		NodeAffinity:  []string{},
	}

	suite.session.Insert(entity.PodCollectionName, &pod)
	defer suite.session.Remove(entity.PodCollectionName, "name", pod.Name)

	retPod := entity.Pod{}
	err := suite.session.FindOne(entity.PodCollectionName, bson.M{"name": pod.Name}, &retPod)
	suite.NoError(err)

	pod, err = FindPodByID(suite.session, retPod.ID)
	suite.NoError(err)

	pod, err = FindPodByID(suite.session, "nonono")
	suite.Error(err)
}
