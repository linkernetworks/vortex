package volume

import (
	"fmt"
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"gopkg.in/mgo.v2/bson"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getPVCInstance(volume *entity.Volume, name string, storageClassName string) *v1.PersistentVolumeClaim {
	capacity, _ := resource.ParseQuantity(volume.Capacity)
	return &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{volume.AccessMode},
			Resources: v1.ResourceRequirements{
				Limits: map[v1.ResourceName]resource.Quantity{
					"storage": capacity,
				},
				Requests: map[v1.ResourceName]resource.Quantity{
					"storage": capacity,
				},
			},
			StorageClassName: &storageClassName,
		},
	}
}

func getStorageClassName(session *mongo.Session, storageName string) (string, error) {
	storage := entity.Storage{}
	err := session.FindOne(entity.StorageCollectionName, bson.M{"name": storageName}, &storage)
	return storage.StorageClassName, err
}

func CreateVolume(sp *serviceprovider.Container, volume *entity.Volume) error {
	session := sp.Mongo.NewSession()
	defer session.Close()
	//fetch the db to get the storageName
	storageName, err := getStorageClassName(session, volume.StorageName)
	if err != nil {
		return err
	}

	name := volume.GetPVCName()
	pvc := getPVCInstance(volume, name, storageName)
	_, err = sp.KubeCtl.CreatePVC(pvc)
	return err
}

func DeleteVolume(sp *serviceprovider.Container, volume *entity.Volume) error {
	//Check the pod
	session := sp.Mongo.NewSession()
	defer session.Close()

	pods := []entity.Pod{}
	fmt.Println("Volume.Name", volume.Name)
	err := session.FindAll(entity.PodCollectionName, bson.M{"volumes.name": volume.Name}, &pods)
	if err != nil {
		return fmt.Errorf("Load the database fail:%v", err)
	}

	for _, pod := range pods {
		//Check the pod's status, report error if at least one pod is running.
		fmt.Println(pod)
	}

	return sp.KubeCtl.DeletePVC(volume.GetPVCName())
}
