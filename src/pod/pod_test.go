package pod

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

type PodTestSuite struct {
	suite.Suite
	sp *serviceprovider.Container
}

func (suite *PodTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	suite.sp = serviceprovider.NewForTesting(cf)
}

func (suite *PodTestSuite) TearDownSuite() {
}

func TestPodSuite(t *testing.T) {
	suite.Run(t, new(PodTestSuite))
}

func (suite *PodTestSuite) TestCreatePod() {
	containers := []entity.Container{
		{
			Name:    namesgenerator.GetRandomName(0),
			Image:   "busybox",
			Command: []string{"sleep", "3600"},
		},
	}

	podName := namesgenerator.GetRandomName(0)
	pod := &entity.Pod{
		ID:         bson.NewObjectId(),
		Name:       podName,
		Containers: containers,
	}

	err := CreatePod(suite.sp, pod)
	suite.NoError(err)

	err = DeletePod(suite.sp, podName)
	suite.NoError(err)
}
