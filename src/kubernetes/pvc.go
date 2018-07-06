package kubernetes

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Get the PVC object by the PVC name
func (kc *KubeCtl) GetPVC(name string) (*corev1.PersistentVolumeClaim, error) {
	return kc.Clientset.CoreV1().PersistentVolumeClaims(kc.Namespace).Get(name, metav1.GetOptions{})
}

//Get all PVCs from the k8s cluster
func (kc *KubeCtl) GetPVCs() ([]*corev1.PersistentVolumeClaim, error) {
	pvcs := []*corev1.PersistentVolumeClaim{}
	pvcsList, err := kc.Clientset.CoreV1().PersistentVolumeClaims(kc.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return pvcs, err
	}
	for _, p := range pvcsList.Items {
		pvcs = append(pvcs, &p)
	}
	return pvcs, nil
}

//Create the PVC by the PVC object
func (kc *KubeCtl) CreatePVC(pvc *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
	return kc.Clientset.CoreV1().PersistentVolumeClaims(kc.Namespace).Create(pvc)
}

//Delete the PVC by the PVC name
func (kc *KubeCtl) DeletePVC(name string) error {
	options := metav1.DeleteOptions{}
	return kc.Clientset.CoreV1().PersistentVolumeClaims(kc.Namespace).Delete(name, &options)
}
