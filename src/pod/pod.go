package pod

import (
	"fmt"
	"strconv"

	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/linkernetworks/vortex/src/utils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"gopkg.in/mgo.v2/bson"
)

var allCapabilities = []corev1.Capability{"NET_ADMIN", "SYS_ADMIN", "NET_RAW"}

// VolumeNamePrefix will set prefix of volumename
const VolumeNamePrefix = "volume-"

// CheckPodParameter will Check Pod's Parameter
func CheckPodParameter(sp *serviceprovider.Container, pod *entity.Pod) error {
	session := sp.Mongo.NewSession()
	defer session.Close()

	//Check the volume
	for _, v := range pod.Volumes {
		count, err := session.Count(entity.VolumeCollectionName, bson.M{"name": v.Name})
		if err != nil {
			return fmt.Errorf("Check the volume name error:%v", err)
		} else if count == 0 {
			return fmt.Errorf("The volume name %s doesn't exist", v.Name)
		}
	}

	//Check the network
	for _, v := range pod.Networks {
		count, err := session.Count(entity.NetworkCollectionName, bson.M{"name": v.Name})
		if err != nil {
			return fmt.Errorf("check the network name error:%v", err)
		} else if count == 0 {
			return fmt.Errorf("the network named %s doesn't exist", v.Name)
		}
	}

	return nil
}

func generateVolume(session *mongo.Session, pod *entity.Pod) ([]corev1.Volume, []corev1.VolumeMount, error) {
	volumes := []corev1.Volume{}
	volumeMounts := []corev1.VolumeMount{}

	for i, v := range pod.Volumes {
		volume := entity.Volume{}
		if err := session.FindOne(entity.VolumeCollectionName, bson.M{"name": v.Name}, &volume); err != nil {
			return nil, nil, fmt.Errorf("Get the volume object error:%v", err)
		}

		vName := fmt.Sprintf("%s-%d", VolumeNamePrefix, i)

		volumes = append(volumes, corev1.Volume{
			Name: vName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: volume.GetPVCName(),
				},
			},
		})

		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      vName,
			MountPath: v.MountPath,
		})
	}

	return volumes, volumeMounts, nil
}

//Get the Intersection of nodes' name
func generateNodeLabels(networks []entity.Network) []string {
	totalNames := [][]string{}
	for _, network := range networks {
		names := []string{}
		for _, node := range network.Nodes {
			names = append(names, node.Name)
		}

		totalNames = append(totalNames, names)
	}

	return utils.Intersections(totalNames)
}

func generateClientCommand(network entity.PodNetwork) (command []string) {
	ip := utils.IPToCIDR(network.IPAddress, network.Netmask)

	command = []string{
		"-s=unix:///tmp/vortex.sock",
		"-b=" + network.BridgeName,
		"-n=" + network.IfName,
		"-i=" + ip,
	}

	if network.VlanTag != nil {
		command = append(command, "-v="+strconv.Itoa((int)(*network.VlanTag)))
	}
	// Support one command with one add route in first version
	if network.Routes != nil {
		if network.Routes[0].Gateway != "" {
			command = append(command, "--net="+network.Routes[0].DstCIDR, "-g="+network.Routes[0].Gateway)
		} else {
			command = append(command, "--net="+network.Routes[0].DstCIDR)
		}
	}
	return
}

func generateInitContainer(networks []entity.PodNetwork) ([]corev1.Container, error) {
	containers := []corev1.Container{}

	for i, v := range networks {
		containers = append(containers, corev1.Container{
			Name:    fmt.Sprintf("init-network-client-%d", i),
			Image:   "sdnvortex/network-controller:v0.4.0",
			Command: []string{"/go/bin/client"},
			Args:    generateClientCommand(v),
			Env: []corev1.EnvVar{
				{
					Name: "POD_NAME",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							FieldPath: "metadata.name",
						},
					},
				},
				{
					Name: "POD_NAMESPACE",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							FieldPath: "metadata.namespace",
						},
					},
				},
				{
					Name: "POD_UUID",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							FieldPath: "metadata.uid",
						},
					},
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "grpc-sock",
					MountPath: "/tmp/",
				},
			},
		})
	}

	return containers, nil
}

//For the network, we will generate two things
//[]string => a list of nodes and it will apply on nodeaffinity
//[]corev1.Container => a list of init container we will apply on pod
func generateNetwork(session *mongo.Session, pod *entity.Pod) ([]string, []corev1.Container, error) {

	networks := []entity.Network{}
	for i, v := range pod.Networks {
		network := entity.Network{}
		if err := session.FindOne(entity.NetworkCollectionName, bson.M{"name": v.Name}, &network); err != nil {
			return nil, nil, err
		}
		networks = append(networks, network)
		pod.Networks[i].BridgeName = network.BridgeName
	}

	nodes := generateNodeLabels(networks)
	containers, err := generateInitContainer(pod.Networks)
	return nodes, containers, err
}

func generateContainerSecurity(pod *entity.Pod) *corev1.SecurityContext {
	if !pod.Capability {
		return &corev1.SecurityContext{}
	}

	privileged := true
	return &corev1.SecurityContext{
		Privileged: &privileged,
		Capabilities: &corev1.Capabilities{
			Add: allCapabilities,
		},
	}
}

func generateAffinity(nodeNames []string) *corev1.Affinity {
	if len(nodeNames) == 0 {
		return &corev1.Affinity{}
	}
	return &corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
				NodeSelectorTerms: []corev1.NodeSelectorTerm{
					{
						MatchExpressions: []corev1.NodeSelectorRequirement{
							{
								Key:      "kubernetes.io/hostname",
								Values:   nodeNames,
								Operator: corev1.NodeSelectorOpIn,
							},
						},
					},
				},
			},
		},
	}
}

// CreatePod will Create Pod
func CreatePod(sp *serviceprovider.Container, pod *entity.Pod) error {
	session := sp.Mongo.NewSession()
	defer session.Close()

	volumes, volumeMounts, err := generateVolume(session, pod)
	if err != nil {
		return err
	}

	nodeNames := []string{}
	initContainers := []corev1.Container{}
	hostNetwork := false
	switch pod.NetworkType {
	case entity.PodHostNetwork:
		hostNetwork = true
	case entity.PodCustomNetwork:
		nodeNames, initContainers, err = generateNetwork(session, pod)
	case entity.PodClusterNetwork:
		//For cluster network, we won't set the nodeAffinity and any netwokr options.
	default:
		err = fmt.Errorf("UnSupported Pod NetworkType %s", pod.NetworkType)
	}
	if err != nil {
		return err
	}

	volumes = append(volumes, corev1.Volume{
		Name: "grpc-sock",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/tmp/vortex",
			},
		},
	})

	var containers []corev1.Container
	securityContext := generateContainerSecurity(pod)
	for _, container := range pod.Containers {
		containers = append(containers, corev1.Container{
			Name:            container.Name,
			Image:           container.Image,
			Command:         container.Command,
			VolumeMounts:    volumeMounts,
			SecurityContext: securityContext,
		})
	}

	p := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   pod.Name,
			Labels: pod.Labels,
		},
		Spec: corev1.PodSpec{
			InitContainers: initContainers,
			Containers:     containers,
			Volumes:        volumes,
			Affinity:       generateAffinity(nodeNames),
			RestartPolicy:  corev1.RestartPolicy(pod.RestartPolicy),
			HostNetwork:    hostNetwork,
		},
	}

	if pod.Namespace == "" {
		pod.Namespace = "default"
	}
	_, err = sp.KubeCtl.CreatePod(&p, pod.Namespace)
	return err
}

// DeletePod will delete pod
func DeletePod(sp *serviceprovider.Container, pod *entity.Pod) error {
	return sp.KubeCtl.DeletePod(pod.Name, pod.Namespace)
}
