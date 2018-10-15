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

type DeploymentTestSuite struct {
	suite.Suite
	sp      *serviceprovider.Container
	session *mongo.Session
}

func (suite *DeploymentTestSuite) SetupSuite() {
	cf := config.MustRead("../../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	suite.sp = sp
	// init session
	suite.session = sp.Mongo.NewSession()
}

func (suite *DeploymentTestSuite) TearDownSuite() {}

func TestDeploymentSuite(t *testing.T) {
	suite.Run(t, new(DeploymentTestSuite))
}

func (suite *DeploymentTestSuite) TestFindDeploymentByName() {
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
	deploy := entity.Deployment{
		Name:         tName,
		Namespace:    namespace,
		Labels:       map[string]string{},
		EnvVars:      map[string]string{},
		Containers:   containers,
		Volumes:      []entity.DeploymentVolume{},
		ConfigMaps:   []entity.DeploymentConfig{},
		Networks:     []entity.DeploymentNetwork{},
		Capability:   true,
		NetworkType:  entity.DeploymentHostNetwork,
		NodeAffinity: []string{},
		Replicas:     1,
	}

	suite.session.Insert(entity.DeploymentCollectionName, &deploy)
	defer suite.session.Remove(entity.DeploymentCollectionName, "name", deploy.Name)

	//load data to check
	retDeployment := entity.Deployment{}
	err := suite.session.FindOne(entity.DeploymentCollectionName, bson.M{"name": deploy.Name}, &retDeployment)
	suite.NoError(err)
	suite.NotEqual("", retDeployment.ID)
	suite.Equal(deploy.Name, retDeployment.Name)
	suite.Equal(len(deploy.Containers), len(retDeployment.Containers))

	deploy, err = FindDeploymentByName(suite.session, deploy.Name)
	suite.NoError(err)

	deploy, err = FindDeploymentByName(suite.session, "nonono")
	suite.Error(err)
}
