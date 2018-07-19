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

func (suite *PodTestSuite) TestCheckPodParameter() {
	volumeName := namesgenerator.GetRandomName(0)
	networkName := namesgenerator.GetRandomName(0)
	pod := &entity.Pod{
		ID:   bson.NewObjectId(),
		Name: namesgenerator.GetRandomName(0),
		Volumes: []entity.PodVolume{
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

	err := CheckPodParameter(suite.sp, pod)
	suite.NoError(err)
}

func (suite *PodTestSuite) TestCheckPodParameterFail() {
	testCases := []struct {
		caseName string
		pod      *entity.Pod
	}{
		{
			"InvalidPodName", &entity.Pod{
				ID:   bson.NewObjectId(),
				Name: "~!@#$%^&*()",
			},
		},
		{
			"InvalidContainerName", &entity.Pod{
				ID:   bson.NewObjectId(),
				Name: namesgenerator.GetRandomName(0),
				Containers: []entity.Container{
					{
						Name:    "~!@#$%^&*()",
						Image:   "busybox",
						Command: []string{"sleep", "3600"},
					},
				},
			},
		},
		{
			"InvalidVolume", &entity.Pod{
				ID:   bson.NewObjectId(),
				Name: namesgenerator.GetRandomName(0),
				Volumes: []entity.PodVolume{
					{Name: namesgenerator.GetRandomName(0)},
				},
			},
		},
		{
			"InvalidNetwork", &entity.Pod{
				ID:   bson.NewObjectId(),
				Name: namesgenerator.GetRandomName(0),
				Networks: []entity.PodNetwork{
					{Name: namesgenerator.GetRandomName(0)},
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.caseName, func(t *testing.T) {
			err := CheckPodParameter(suite.sp, tc.pod)
			suite.Error(err)
		})
	}
}

func (suite *PodTestSuite) TestGenerateVolume() {
	volumeName := namesgenerator.GetRandomName(0)
	pod := &entity.Pod{
		ID: bson.NewObjectId(),
		Volumes: []entity.PodVolume{
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

	volumes, volumeMounts, err := generateVolume(session, pod)
	suite.NotNil(volumes)
	suite.NotNil(volumeMounts)
	suite.NoError(err)
}

func (suite *PodTestSuite) TestGenerateVolumeFail() {
	volumeName := namesgenerator.GetRandomName(0)
	pod := &entity.Pod{
		ID: bson.NewObjectId(),
		Volumes: []entity.PodVolume{
			{Name: volumeName},
		},
	}

	session := suite.sp.Mongo.NewSession()
	defer session.Close()
	volumes, volumeMounts, err := generateVolume(session, pod)
	suite.Nil(volumes)
	suite.Nil(volumeMounts)
	suite.Error(err)
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

	err = DeletePod(suite.sp, pod)
	suite.NoError(err)
}

func (suite *PodTestSuite) TestCreatePodFailWithoutVolume() {
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
		Volumes: []entity.PodVolume{
			{Name: namesgenerator.GetRandomName(0)},
		},
	}

	err := CreatePod(suite.sp, pod)
	suite.Error(err)
}

func (suite *PodTestSuite) TestCreatePodFailWithoutNetwork() {
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
		Networks: []entity.PodNetwork{
			{
				Name: namesgenerator.GetRandomName(0),
			},
		},
	}

	err := CreatePod(suite.sp, pod)
	suite.Error(err)
}

func (suite *PodTestSuite) TestGenerateNodeLabels() {
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

func (suite *PodTestSuite) TestGenerateClientCommand() {
	bName := namesgenerator.GetRandomName(0)
	ifName := namesgenerator.GetRandomName(0)
	podNetwork := entity.PodNetwork{
		IfName:     ifName,
		IPAddress:  "1.2.3.4",
		Netmask:    "255.255.255.0",
		BridgeName: bName,
	}
	command := generateClientCommand(podNetwork)
	ans := []string{"-s=unix:///tmp/vortex.sock", "-b=" + bName, "-n=" + ifName, "-i=1.2.3.4/24"}
	suite.Equal(ans, command)

	var vlanTag int32
	vlanTag = 123
	podNetwork.VlanTag = &vlanTag
	command = generateClientCommand(podNetwork)
	ans = []string{"-s=unix:///tmp/vortex.sock", "-b=" + bName, "-n=" + ifName, "-i=1.2.3.4/24", "-v=123"}
	suite.Equal(ans, command)

}

func (suite *PodTestSuite) TestGenerateNetwork() {
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

	podName := namesgenerator.GetRandomName(0)
	ifName := namesgenerator.GetRandomName(0)
	pod := &entity.Pod{
		ID:   bson.NewObjectId(),
		Name: podName,
		Networks: []entity.PodNetwork{
			{
				Name:      networkName,
				IfName:    ifName,
				IPAddress: "1.2.3.4",
				Netmask:   "255.255.255.0",
			},
		},
	}

	nodes, containers, err := generateNetwork(session, pod)
	suite.NoError(err)
	suite.Equal(1, len(containers))
	suite.Equal(0, len(nodes))
}

func (suite *PodTestSuite) TestGenerateNetworkFail() {
	networkName := namesgenerator.GetRandomName(0)
	podName := namesgenerator.GetRandomName(0)
	ifName := namesgenerator.GetRandomName(0)
	pod := &entity.Pod{
		ID:   bson.NewObjectId(),
		Name: podName,
		Networks: []entity.PodNetwork{
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

	nodes, containers, err := generateNetwork(session, pod)
	suite.Error(err)
	suite.Nil(nodes)
	suite.Nil(containers)
}

func (suite *PodTestSuite) TestGenerateAffinity() {
	affinity := generateAffinity([]string{})
	suite.Nil(affinity.NodeAffinity)
	affinity = generateAffinity([]string{"123"})
	suite.NotNil(affinity.NodeAffinity)

}
