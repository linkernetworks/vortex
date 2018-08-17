package deploy

import (
	"math/rand"
	"testing"
	"time"

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
	sp *serviceprovider.Container
}

func (suite *DeploymentTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	suite.sp = serviceprovider.NewForTesting(cf)
}

func (suite *DeploymentTestSuite) TearDownSuite() {
}

func TestDeploymentSuite(t *testing.T) {
	suite.Run(t, new(DeploymentTestSuite))
}

func (suite *DeploymentTestSuite) TestCheckDeploymentParameter() {
	volumeName := namesgenerator.GetRandomName(0)
	networkName := namesgenerator.GetRandomName(0)
	deploy := &entity.Deployment{
		ID:   bson.NewObjectId(),
		Name: namesgenerator.GetRandomName(0),
		Volumes: []entity.DeploymentVolume{
			{Name: volumeName},
		},
	}

	session := suite.sp.Mongo.NewSession()
	defer session.Close()

	volume := entity.Volume{
		ID:   bson.NewObjectId(),
		Name: volumeName,
	}
	session.Insert(entity.VolumeCollectionName, volume)
	defer session.Remove(entity.VolumeCollectionName, "name", volume.Name)

	network := entity.Network{
		ID:   bson.NewObjectId(),
		Name: networkName,
	}
	session.Insert(entity.NetworkCollectionName, network)
	defer session.Remove(entity.NetworkCollectionName, "name", network.Name)

	err := CheckDeploymentParameter(suite.sp, deploy)
	suite.NoError(err)
}

func (suite *DeploymentTestSuite) TestCheckDeploymentParameterFail() {
	testCases := []struct {
		caseName string
		deploy   *entity.Deployment
	}{
		{
			"InvalidVolume", &entity.Deployment{
				ID:   bson.NewObjectId(),
				Name: namesgenerator.GetRandomName(0),
				Volumes: []entity.DeploymentVolume{
					{Name: namesgenerator.GetRandomName(0)},
				},
			},
		},
		{
			"InvalidNetwork", &entity.Deployment{
				ID:   bson.NewObjectId(),
				Name: namesgenerator.GetRandomName(0),
				Networks: []entity.DeploymentNetwork{
					{Name: namesgenerator.GetRandomName(0)},
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.caseName, func(t *testing.T) {
			err := CheckDeploymentParameter(suite.sp, tc.deploy)
			suite.Error(err)
		})
	}
}

func (suite *DeploymentTestSuite) TestGenerateVolume() {
	volumeName := namesgenerator.GetRandomName(0)
	deploy := &entity.Deployment{
		ID: bson.NewObjectId(),
		Volumes: []entity.DeploymentVolume{
			{Name: volumeName},
		},
	}

	session := suite.sp.Mongo.NewSession()
	defer session.Close()

	volume := entity.Volume{
		ID:   bson.NewObjectId(),
		Name: volumeName,
	}
	session.Insert(entity.VolumeCollectionName, volume)
	defer session.Remove(entity.VolumeCollectionName, "name", volume.Name)

	volumes, volumeMounts, err := generateVolume(session, deploy)
	suite.NotNil(volumes)
	suite.NotNil(volumeMounts)
	suite.NoError(err)
}

func (suite *DeploymentTestSuite) TestGenerateVolumeFail() {
	volumeName := namesgenerator.GetRandomName(0)
	deploy := &entity.Deployment{
		ID: bson.NewObjectId(),
		Volumes: []entity.DeploymentVolume{
			{Name: volumeName},
		},
	}

	session := suite.sp.Mongo.NewSession()
	defer session.Close()
	volumes, volumeMounts, err := generateVolume(session, deploy)
	suite.Nil(volumes)
	suite.Nil(volumeMounts)
	suite.Error(err)
}

func (suite *DeploymentTestSuite) TestCreateDeployment() {
	containers := []entity.Container{
		{
			Name:    namesgenerator.GetRandomName(0),
			Image:   "busybox",
			Command: []string{"sleep", "3600"},
		},
	}

	deployName := namesgenerator.GetRandomName(0)
	deploy := &entity.Deployment{
		ID:          bson.NewObjectId(),
		Name:        deployName,
		Containers:  containers,
		NetworkType: entity.DeploymentHostNetwork,
		EnvVars: map[string]string{
			"MY_IP": "1.2.3.4",
		},
	}

	err := CreateDeployment(suite.sp, deploy)
	suite.NoError(err)

	err = DeleteDeployment(suite.sp, deploy)
	suite.NoError(err)
}

func (suite *DeploymentTestSuite) TestCreateDeploymentFailWithoutVolume() {
	containers := []entity.Container{
		{
			Name:    namesgenerator.GetRandomName(0),
			Image:   "busybox",
			Command: []string{"sleep", "3600"},
		},
	}

	deployName := namesgenerator.GetRandomName(0)
	deploy := &entity.Deployment{
		ID:         bson.NewObjectId(),
		Name:       deployName,
		Containers: containers,
		Volumes: []entity.DeploymentVolume{
			{Name: namesgenerator.GetRandomName(0)},
		},
	}

	err := CreateDeployment(suite.sp, deploy)
	suite.Error(err)
}

func (suite *DeploymentTestSuite) TestCreateDeploymentFailWithoutNetwork() {
	containers := []entity.Container{
		{
			Name:    namesgenerator.GetRandomName(0),
			Image:   "busybox",
			Command: []string{"sleep", "3600"},
		},
	}

	deployName := namesgenerator.GetRandomName(0)
	deploy := &entity.Deployment{
		ID:         bson.NewObjectId(),
		Name:       deployName,
		Containers: containers,
		Networks: []entity.DeploymentNetwork{
			{
				Name: namesgenerator.GetRandomName(0),
			},
		},
	}

	err := CreateDeployment(suite.sp, deploy)
	suite.Error(err)
}

func (suite *DeploymentTestSuite) TestGenerateNodeLabels() {
	networks := []entity.Network{
		{
			Nodes: []entity.Node{
				{Name: "node1"},
				{Name: "node4"},
				{Name: "node5"},
			},
		},
		{
			Nodes: []entity.Node{
				{Name: "node2"},
				{Name: "node3"},
				{Name: "node4"},
				{Name: "node5"},
				{Name: "node6"},
			},
		},
		{
			Nodes: []entity.Node{
				{Name: "node1"},
				{Name: "node2"},
				{Name: "node3"},
				{Name: "node4"},
				{Name: "node5"},
			},
		},
	}

	names := generateNodeLabels(networks)
	suite.Equal(2, len(names))
	suite.Equal([]string{"node4", "node5"}, names)
}

func (suite *DeploymentTestSuite) TestGenerateClientCommand() {
	bName := namesgenerator.GetRandomName(0)
	ifName := namesgenerator.GetRandomName(0)
	deployNetwork := entity.DeploymentNetwork{
		IfName:     ifName,
		IPAddress:  "1.2.3.4",
		Netmask:    "255.255.255.0",
		BridgeName: bName,
	}
	command := generateClientCommand(deployNetwork)
	ans := []string{"-s=unix:///tmp/vortex.sock", "-b=" + bName, "-n=" + ifName, "-i=1.2.3.4/24"}
	suite.Equal(ans, command)

	var vlanTag int32
	vlanTag = 123
	deployNetwork.VlanTag = &vlanTag
	command = generateClientCommand(deployNetwork)
	ans = []string{"-s=unix:///tmp/vortex.sock", "-b=" + bName, "-n=" + ifName, "-i=1.2.3.4/24", "-v=123"}
	suite.Equal(ans, command)

}

func (suite *DeploymentTestSuite) TestGenerateNetwork() {
	networkName := namesgenerator.GetRandomName(0)
	bName := namesgenerator.GetRandomName(0)
	network := entity.Network{
		ID:         bson.NewObjectId(),
		Name:       networkName,
		BridgeName: bName,
	}
	session := suite.sp.Mongo.NewSession()
	defer session.Close()

	session.Insert(entity.NetworkCollectionName, network)
	defer session.Remove(entity.NetworkCollectionName, "name", network.Name)

	deployName := namesgenerator.GetRandomName(0)
	ifName := namesgenerator.GetRandomName(0)
	deploy := &entity.Deployment{
		ID:   bson.NewObjectId(),
		Name: deployName,
		Networks: []entity.DeploymentNetwork{
			{
				Name:      networkName,
				IfName:    ifName,
				IPAddress: "1.2.3.4",
				Netmask:   "255.255.255.0",
			},
		},
	}

	nodes, containers, err := generateNetwork(session, deploy)
	suite.NoError(err)
	suite.Equal(1, len(containers))
	suite.Equal(0, len(nodes))
}

func (suite *DeploymentTestSuite) TestGenerateNetworkFail() {
	networkName := namesgenerator.GetRandomName(0)
	deployName := namesgenerator.GetRandomName(0)
	ifName := namesgenerator.GetRandomName(0)
	deploy := &entity.Deployment{
		ID:   bson.NewObjectId(),
		Name: deployName,
		Networks: []entity.DeploymentNetwork{
			{
				Name:      networkName,
				IfName:    ifName,
				IPAddress: "1.2.3.4",
				Netmask:   "255.255.255.0",
			},
		},
	}

	session := suite.sp.Mongo.NewSession()
	defer session.Close()

	nodes, containers, err := generateNetwork(session, deploy)
	suite.Error(err)
	suite.Nil(nodes)
	suite.Nil(containers)
}

func (suite *DeploymentTestSuite) TestGenerateAffinity() {
	affinity := generateAffinity([]string{})
	suite.Nil(affinity.NodeAffinity)
	affinity = generateAffinity([]string{"123"})
	suite.NotNil(affinity.NodeAffinity)
}

func (suite *DeploymentTestSuite) TestGenerateContainerSecurityContext() {
	deploy := &entity.Deployment{}
	security := generateContainerSecurity(deploy)
	suite.Nil(security.Privileged)
	suite.Nil(security.Capabilities)

	deploy.Capability = true
	security = generateContainerSecurity(deploy)
	suite.NotNil(security.Privileged)
	suite.NotNil(security.Capabilities)
}

func (suite *DeploymentTestSuite) TestCreateDeploymentWithNetworkTypes() {

	networkName := namesgenerator.GetRandomName(0)
	bName := namesgenerator.GetRandomName(0)
	network := entity.Network{
		ID:         bson.NewObjectId(),
		Name:       networkName,
		BridgeName: bName,
		Nodes: []entity.Node{
			{Name: "node1"},
			{Name: "node2"},
			{Name: "node3"},
		},
	}
	session := suite.sp.Mongo.NewSession()
	defer session.Close()

	session.Insert(entity.NetworkCollectionName, network)
	defer session.Remove(entity.NetworkCollectionName, "name", network.Name)

	//For each case, we need to check the
	//HostNetwork object
	//NodeAffinity object
	//Create Success
	testCases := []struct {
		caseName      string
		deploy        *entity.Deployment
		rHostNetwork  bool     //result of hostNetwork
		rNodeAffinity []string //result of rNodeAffinity
	}{
		{
			"hostNetwork",
			&entity.Deployment{
				ID:           bson.NewObjectId(),
				Name:         namesgenerator.GetRandomName(0),
				Containers:   []entity.Container{},
				NetworkType:  entity.DeploymentHostNetwork,
				NodeAffinity: []string{"node1", "node2"},
			},
			true,
			[]string{"node1", "node2"},
		},
		{
			"clusterNetwork",
			&entity.Deployment{
				ID:           bson.NewObjectId(),
				Name:         namesgenerator.GetRandomName(0),
				Containers:   []entity.Container{},
				NetworkType:  entity.DeploymentClusterNetwork,
				NodeAffinity: []string{"node1", "node2"},
			},
			false,
			[]string{"node1", "node2"},
		},
		{
			"customNetwork",
			&entity.Deployment{
				ID:          bson.NewObjectId(),
				Name:        namesgenerator.GetRandomName(0),
				Containers:  []entity.Container{},
				NetworkType: entity.DeploymentCustomNetwork,
				Networks: []entity.DeploymentNetwork{
					{
						Name: networkName,
					},
				},
				NodeAffinity: []string{"node1", "node5"},
			},
			false,
			[]string{"node1"},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.caseName, func(t *testing.T) {
			err := CreateDeployment(suite.sp, tc.deploy)
			suite.NoError(err)

			deploy, err := suite.sp.KubeCtl.GetDeployment(tc.deploy.Name, tc.deploy.Namespace)
			suite.NotNil(deploy)
			suite.NoError(err)

			suite.Equal(tc.rHostNetwork, deploy.Spec.Template.Spec.HostNetwork)
			suite.Equal(tc.rNodeAffinity, deploy.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Values)

			err = DeleteDeployment(suite.sp, tc.deploy)
			suite.NoError(err)

		})
	}
}
