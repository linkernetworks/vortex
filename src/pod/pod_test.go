package pod

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
		ID:          bson.NewObjectId(),
		Name:        podName,
		Containers:  containers,
		NetworkType: entity.PodHostNetwork,
		EnvVars: map[string]string{
			"MY_IP": "1.2.3.4",
		},
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
		Name:      "my-net",
		IfName:    ifName,
		IPAddress: "1.2.3.4",
		Netmask:   "255.255.255.0",
		RoutesGw: []entity.PodRouteGw{
			{
				DstCIDR: "192.168.2.0/24",
				Gateway: "192.168.2.254",
			},
		},
		RoutesIntf: []entity.PodRouteIntf{
			{
				DstCIDR: "192.168.3.0/24",
			},
		},
		BridgeName: bName,
	}
	command := generateClientCommand(podNetwork)
	ans := []string{
		"--server=unix:///tmp/vortex.sock",
		"--bridge=" + bName,
		"--nic=" + ifName,
		"--ip=1.2.3.4/24",
		"--route-gw=192.168.2.0/24,192.168.2.254",
		"--route-intf=192.168.3.0/24",
	}
	suite.Equal(ans, command)

	var vlanTag int32
	vlanTag = 123
	podNetwork.VlanTag = &vlanTag
	command = generateClientCommand(podNetwork)
	ans = []string{
		"--server=unix:///tmp/vortex.sock",
		"--bridge=" + bName,
		"--nic=" + ifName,
		"--ip=1.2.3.4/24",
		"--vlan=123",
		"--route-gw=192.168.2.0/24,192.168.2.254",
		"--route-intf=192.168.3.0/24",
	}
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

func (suite *PodTestSuite) TestGenerateContainerSecurityContext() {
	pod := &entity.Pod{}
	security := generateContainerSecurity(pod)
	suite.Nil(security.Privileged)
	suite.Nil(security.Capabilities)

	pod.Capability = true
	security = generateContainerSecurity(pod)
	suite.NotNil(security.Privileged)
	suite.NotNil(security.Capabilities)
}

func (suite *PodTestSuite) TestCreatePodWithNetworkTypes() {

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
		pod           *entity.Pod
		rHostNetwork  bool     //result of hostNetwork
		rNodeAffinity []string //result of rNodeAffinity
	}{
		{
			"hostNetwork",
			&entity.Pod{
				ID:           bson.NewObjectId(),
				Name:         namesgenerator.GetRandomName(0),
				Containers:   []entity.Container{},
				NetworkType:  entity.PodHostNetwork,
				NodeAffinity: []string{"node1", "node2"},
			},
			true,
			[]string{"node1", "node2"},
		},
		{
			"clusterNetwork",
			&entity.Pod{
				ID:           bson.NewObjectId(),
				Name:         namesgenerator.GetRandomName(0),
				Containers:   []entity.Container{},
				NetworkType:  entity.PodClusterNetwork,
				NodeAffinity: []string{"node1", "node2"},
			},
			false,
			[]string{"node1", "node2"},
		},
		{
			"customNetwork",
			&entity.Pod{
				ID:          bson.NewObjectId(),
				Name:        namesgenerator.GetRandomName(0),
				Containers:  []entity.Container{},
				NetworkType: entity.PodCustomNetwork,
				Networks: []entity.PodNetwork{
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
			err := CreatePod(suite.sp, tc.pod)
			suite.NoError(err)

			pod, err := suite.sp.KubeCtl.GetPod(tc.pod.Name, tc.pod.Namespace)
			suite.NotNil(pod)
			suite.NoError(err)

			suite.Equal(tc.rHostNetwork, pod.Spec.HostNetwork)
			suite.Equal(tc.rNodeAffinity, pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Values)

			err = DeletePod(suite.sp, tc.pod)
			suite.NoError(err)

		})
	}
}
