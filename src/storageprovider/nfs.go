package storageprovider

import (
	"fmt"
	"net"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"gopkg.in/mgo.v2/bson"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	NFSProvisionerPrefix  string = "nfs-provisioner-"
	NFSStorageClassPrefix string = "nfs-storageclass-"
)

type NFSStorageProvider struct {
	entity.Storage
}

func (nfs NFSStorageProvider) ValidateBeforeCreating(sp *serviceprovider.Container, storage *entity.Storage) error {
	ip := net.ParseIP(storage.IP)
	if len(ip) == 0 {
		return fmt.Errorf("Invalid IP address %s\n", storage.IP)
	}

	path := storage.PATH
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
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
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
					ServiceAccountName: "nfs-client-provisioner",
					Containers: []v1.Container{
						{
							Name:            name,
							Image:           "quay.io/external_storage/nfs-client-provisioner:latest",
							ImagePullPolicy: v1.PullIfNotPresent,
							Env: []v1.EnvVar{
								{Name: "PROVISIONER_NAME", Value: name},
								{Name: "NFS_SERVER", Value: storage.IP},
								{Name: "NFS_PATH", Value: storage.PATH},
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
									Server: storage.IP,
									Path:   storage.PATH,
								},
							},
						},
					},
				},
			},
		},
	}
}

func getStorageClass(name string, provisioner string, storage *entity.Storage) *storagev1.StorageClass {
	return &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Provisioner: provisioner,
	}
}

func (nfs NFSStorageProvider) CreateStorage(sp *serviceprovider.Container, storage *entity.Storage) error {
	namespace := "default"
	name := NFSProvisionerPrefix + storage.ID.Hex()
	storageClassName := NFSStorageClassPrefix + storage.ID.Hex()
	//Create deployment
	deployment := getDeployment(name, storage)
	//Create storageClass
	storageClass := getStorageClass(storageClassName, name, storage)
	storage.StorageClassName = storageClassName
	if _, err := sp.KubeCtl.CreateDeployment(deployment, namespace); err != nil {
		return err
	}
	_, err := sp.KubeCtl.CreateStorageClass(storageClass)
	return err
}

func (nfs NFSStorageProvider) DeleteStorage(sp *serviceprovider.Container, storage *entity.Storage) error {
	namespace := "default"
	deployName := NFSProvisionerPrefix + storage.ID.Hex()
	storageName := NFSStorageClassPrefix + storage.ID.Hex()

	//If the storage is used by some volume, we can't delete it.
	q := bson.M{"storageName": storage.Name}
	session := sp.Mongo.NewSession()
	defer session.Close()

	count, err := session.Count(entity.VolumeCollectionName, q)
	if err != nil {
		return err
	} else if count > 0 {
		return fmt.Errorf("The StorageName %s can't be deleted, since there're some volume still use it", storage.Name)
	}

	//Delete StorageClass
	if err := sp.KubeCtl.DeleteStorageClass(storageName); err != nil {
		return err
	}
	//Delete Deployment
	return sp.KubeCtl.DeleteDeployment(deployName, namespace)
}
