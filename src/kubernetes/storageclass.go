package kubernetes

import (
	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Get the external IP address of node
func (kc *KubeCtl) CreateStorageClass(storageClass *v1.StorageClass) (*v1.StorageClass, error) {
	return kc.Clientset.StorageV1().StorageClasses().Create(storageClass)
}

func (kc *KubeCtl) GetStorageClass(name string) (*v1.StorageClass, error) {
	return kc.Clientset.StorageV1().StorageClasses().Get(name, metav1.GetOptions{})
}

func (kc *KubeCtl) DeleteStorageClass(name string) error {
	return kc.Clientset.StorageV1().StorageClasses().Delete(name, &metav1.DeleteOptions{})
}
