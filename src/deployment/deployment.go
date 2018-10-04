package deployment

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/linkernetworks/vortex/src/utils"

	appsv1 "k8s.io/api/apps/v1"
	v2beta1 "k8s.io/api/autoscaling/v2beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"gopkg.in/mgo.v2/bson"
)

var allCapabilities = []corev1.Capability{"NET_ADMIN", "SYS_ADMIN", "NET_RAW"}

// VolumeNamePrefix will set prefix of volumename
const VolumeNamePrefix = "volume"

// ConfigMapNamePrefix will set prefix of volumename
const ConfigMapNamePrefix = "configmap"

// DefaultLabel is the label we used for our deploying application/deployment/pods
const DefaultLabel = "vortex"

// NotificationEmailAccount is the label we do notify to user email account
const NotificationEmailAccount = "email_account"

// NotificationEmailDomain is the label we do notify to user email domain
const NotificationEmailDomain = "email_domain"

// CheckDeploymentParameter will Check Deployment's Parameter
func CheckDeploymentParameter(sp *serviceprovider.Container, deploy *entity.Deployment) error {
	session := sp.Mongo.NewSession()
	defer session.Close()

	//Check the volume
	for _, v := range deploy.Volumes {
		count, err := session.Count(entity.VolumeCollectionName, bson.M{"name": v.Name})
		if err != nil {
			return fmt.Errorf("Check the volume name error:%v", err)
		} else if count == 0 {
			return fmt.Errorf("The volume name %s doesn't exist", v.Name)
		}
	}

	//Check the network
	for _, v := range deploy.Networks {
		count, err := session.Count(entity.NetworkCollectionName, bson.M{"name": v.Name})
		if err != nil {
			return fmt.Errorf("check the network name error:%v", err)
		} else if count == 0 {
			return fmt.Errorf("the network named %s doesn't exist", v.Name)
		}
	}

	return nil
}

func generateVolume(session *mongo.Session, deploy *entity.Deployment) ([]corev1.Volume, []corev1.VolumeMount, error) {
	volumes := []corev1.Volume{}
	volumeMounts := []corev1.VolumeMount{}

	for i, v := range deploy.Volumes {
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

func generateConfigMap(deploy *entity.Deployment) ([]corev1.Volume, []corev1.VolumeMount, error) {
	volumes := []corev1.Volume{}
	volumeMounts := []corev1.VolumeMount{}

	for _, v := range deploy.ConfigMaps {

		// TODO: check whether this configMap exist

		vName := fmt.Sprintf("%s-%s", ConfigMapNamePrefix, v.Name)

		volumes = append(volumes, corev1.Volume{
			Name: vName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: v.Name,
					},
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

func generateClientCommand(network entity.DeploymentNetwork) (command []string) {
	ip := utils.IPToCIDR(network.IPAddress, network.Netmask)

	command = []string{
		"--server=unix:///tmp/vortex.sock",
		"--bridge=" + network.BridgeName,
		"--nic=" + network.IfName,
		"--ip=" + ip,
	}

	if network.VlanTag != nil {
		command = append(command, "--vlan="+strconv.Itoa((int)(*network.VlanTag)))
	}
	if len(network.RoutesGw) != 0 {
		for _, netroute := range network.RoutesGw {
			command = append(command, "--route-gw="+netroute.DstCIDR+","+netroute.Gateway)
		}
	}
	if len(network.RoutesIntf) != 0 {
		for _, netroute := range network.RoutesIntf {
			command = append(command, "--route-intf="+netroute.DstCIDR)
		}
	}
	return
}

func generateInitContainer(networks []entity.DeploymentNetwork) ([]corev1.Container, error) {
	containers := []corev1.Container{}

	ethtools := []string{}
	for i, v := range networks {
		containers = append(containers, corev1.Container{
			Name:    fmt.Sprintf("init-network-client-%d", i),
			Image:   "sdnvortex/network-controller:v0.4.8",
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

		if strings.HasPrefix(v.BridgeName, "netdev") {
			ethtools = append(ethtools, v.IfName)
		}
	}

	if len(ethtools) != 0 {
		privileged := true
		containers = append(containers, corev1.Container{
			Name:    "init-ethtool",
			Image:   "sdnvortex/ethtool:latest",
			Command: []string{"/usr/bin/ethtool.sh"},
			Args:    ethtools,
			SecurityContext: &corev1.SecurityContext{
				Privileged: &privileged,
			},
		})
	}

	return containers, nil
}

//For the network, we will generate two things
//[]string => a list of nodes and it will apply on nodeaffinity
//[]corev1.Container => a list of init container we will apply on deploy
func generateNetwork(session *mongo.Session, deploy *entity.Deployment) ([]string, []corev1.Container, error) {
	networks := []entity.Network{}
	for i, v := range deploy.Networks {
		network := entity.Network{}
		if err := session.FindOne(entity.NetworkCollectionName, bson.M{"name": v.Name}, &network); err != nil {
			return nil, nil, err
		}
		networks = append(networks, network)
		deploy.Networks[i].BridgeName = network.BridgeName
	}

	nodes := generateNodeLabels(networks)
	containers, err := generateInitContainer(deploy.Networks)
	return nodes, containers, err
}

func generateContainerSecurity(deploy *entity.Deployment) *corev1.SecurityContext {
	if !deploy.Capability {
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

func generateEnvVars(deploy *entity.Deployment) []corev1.EnvVar {
	envVars := []corev1.EnvVar{}

	for k, v := range deploy.EnvVars {
		envVars = append(envVars, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}
	return envVars
}

// CreateDeployment will Create Deployment
func CreateDeployment(sp *serviceprovider.Container, deploy *entity.Deployment) error {
	session := sp.Mongo.NewSession()
	defer session.Close()

	volumes, volumeMounts, err := generateVolume(session, deploy)
	if err != nil {
		return err
	}

	configMaps, configMapMounts, err := generateConfigMap(deploy)
	if err != nil {
		return err
	}

	nodeAffinity := deploy.NodeAffinity
	initContainers := []corev1.Container{}
	hostNetwork := false
	switch deploy.NetworkType {
	case entity.DeploymentHostNetwork:
		hostNetwork = true
	case entity.DeploymentCustomNetwork:
		var tmp []string
		tmp, initContainers, err = generateNetwork(session, deploy)
		if len(tmp) != 0 {
			nodeAffinity = utils.Intersection(nodeAffinity, tmp)
		}
	case entity.DeploymentClusterNetwork:
		// For cluster network, we won't set the nodeAffinity and any network options.
	default:
		err = fmt.Errorf("Unsupported Deployment NetworkType %s", deploy.NetworkType)
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

	volumes = append(volumes, configMaps...)
	volumeMounts = append(volumeMounts, configMapMounts...)

	var containers []corev1.Container
	var c corev1.Container

	securityContext := generateContainerSecurity(deploy)
	envVars := generateEnvVars(deploy)

	for _, deployContainer := range deploy.Containers {
		c = corev1.Container{
			Name:            deployContainer.Name,
			Image:           deployContainer.Image,
			Command:         deployContainer.Command,
			VolumeMounts:    volumeMounts,
			SecurityContext: securityContext,
			Env:             envVars,
		}
		if deployContainer.ResourceRequestCPU != 0 {
			c.Resources = corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					"cpu": resource.MustParse(strconv.Itoa(deployContainer.ResourceRequestCPU) + "m"),
				},
			}
		} else if deployContainer.ResourceRequestMemory != 0 {
			c.Resources = corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					"memory": resource.MustParse(strconv.Itoa(deployContainer.ResourceRequestMemory) + "Mi"),
				},
			}
		} else if deployContainer.ResourceRequestMemory != 0 && deployContainer.ResourceRequestCPU != 0 {
			c.Resources = corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					"cpu":    resource.MustParse(strconv.Itoa(deployContainer.ResourceRequestCPU) + "m"),
					"memory": resource.MustParse(strconv.Itoa(deployContainer.ResourceRequestMemory) + "Mi"),
				},
			}
		}
		containers = append(containers, c)
	}

	p := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   deploy.Name,
			Labels: deploy.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					// vortex default label
					DefaultLabel: deploy.Name,
				},
			},
			Replicas: &deploy.Replicas,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						// vortex default label
						DefaultLabel: deploy.Name,
					},
				},
				Spec: corev1.PodSpec{
					InitContainers: initContainers,
					Containers:     containers,
					Volumes:        volumes,
					Affinity:       generateAffinity(nodeAffinity),
					RestartPolicy:  corev1.RestartPolicyAlways,
					HostNetwork:    hostNetwork,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{Name: "dockerhub-token"},
					},
				},
			},
		},
	}

	// pass the same labels to pod
	for k, v := range deploy.Labels {
		p.Spec.Template.ObjectMeta.Labels[k] = v
	}

	_, err = sp.KubeCtl.CreateDeployment(&p, deploy.Namespace)
	return err
}

// DeleteDeployment will delete a deployment
func DeleteDeployment(sp *serviceprovider.Container, deploy *entity.Deployment) error {
	return sp.KubeCtl.DeleteDeployment(deploy.Name, deploy.Namespace)
}

// CreateAutoscaler will create a autoscaler
func CreateAutoscaler(sp *serviceprovider.Container, autoscalerInfo entity.AutoscalerInfo) error {
	autoscaler := v2beta1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			// use deployment name to name autoscaler's name
			Name:      autoscalerInfo.ScaleTargetRefName,
			Namespace: autoscalerInfo.Namespace,
		},
		Spec: v2beta1.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: v2beta1.CrossVersionObjectReference{
				APIVersion: "extensions/v1beta",
				Kind:       "Deployment",
				Name:       autoscalerInfo.ScaleTargetRefName,
			},
			MinReplicas: &autoscalerInfo.MinReplicas,
			MaxReplicas: autoscalerInfo.MaxReplicas,
			Metrics: []v2beta1.MetricSpec{
				{
					Type: v2beta1.ResourceMetricSourceType,
					Resource: &v2beta1.ResourceMetricSource{
						Name:                     autoscalerInfo.ResourceName,
						TargetAverageUtilization: &autoscalerInfo.TargetAverageUtilization,
					},
				},
			},
		},
	}
	_, err := sp.KubeCtl.CreateAutoscaler(&autoscaler, autoscalerInfo.Namespace)
	return err
}

// DeleteAutoscaler will delete a autoscaler
func DeleteAutoscaler(sp *serviceprovider.Container, autoscalerInfo entity.AutoscalerInfo) error {
	return sp.KubeCtl.DeleteAutoscaler(autoscalerInfo.ScaleTargetRefName, autoscalerInfo.Namespace)
}
