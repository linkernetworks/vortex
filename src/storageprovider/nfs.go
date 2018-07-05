package storageprovider

import (
	"fmt"
	"net"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//	"gopkg.in/mgo.v2/bson"
)

const NFS_PROVISIONER_PREFIX = "nfs-provisioner-"
const NFS_STORAGECLASS_PREFIX = "nfs-storageclass-"

type NFSStorageProvider struct {
	entity.NFSStorage
}

func (nfs NFSStorageProvider) ValidateBeforeCreating(sp *serviceprovider.Container, storage *entity.Storage) error {
	ip := net.ParseIP(storage.NFS.IP)
	if len(ip) == 0 {
		return fmt.Errorf("Invalid IP address %s\n", storage.NFS.IP)
	}

	path := storage.NFS.PATH
	if path == "" || path[0] != '/' {
		return fmt.Errorf("Invalid NFS export path %s\n", path)
	}

	return nil
}

func getDeployment(name string, storage *entity.Storage) *appsv1.Deployment {
	var replicas int32
	replicas = 1
	volumeName := "nfs-client-root"
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            name,
							Image:           "quay.io/kubernetes_incubator/nfs-provisioner:latest",
							ImagePullPolicy: v1.PullIfNotPresent,
							Env: []v1.EnvVar{
								{Name: "PROVISIONER_NAME", Value: name},
								{Name: "NFS_SERVER", Value: storage.NFS.IP},
								{Name: "NFS_PATH", Value: storage.NFS.PATH},
							},
							VolumeMounts: []v1.VolumeMount{
								{Name: volumeName, MountPath: "/persistentvolumes"},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: volumeName,
							VolumeSource: v1.VolumeSource{
								NFS: &v1.NFSVolumeSource{
									Server: storage.NFS.IP,
									Path:   storage.NFS.PATH,
								},
							},
						},
					},
				},
			},
		},
	}

}

func (nfs NFSStorageProvider) CreateStorage(sp *serviceprovider.Container, storage *entity.Storage) error {
	name := NFS_PROVISIONER_PREFIX + storage.ID.Hex()

	//Create deployment
	deployment := getDeployment(name, storage)
	//Create storageClass
	_, err := sp.KubeCtl.CreateDeployment(deployment)
	return err
}

func (nfs NFSStorageProvider) DeleteStorage(sp *serviceprovider.Container, storage *entity.Storage) error {
	name := NFS_PROVISIONER_PREFIX + storage.ID.Hex()
	//Delete StorageClass

	//Delete Deployment
	return sp.KubeCtl.DeleteDeployment(name)
}
