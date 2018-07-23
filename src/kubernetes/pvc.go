package kubernetes

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPVC will get the PVC object by the PVC name
func (kc *KubeCtl) GetPVC(name string, namespace string) (*corev1.PersistentVolumeClaim, error) {
	return kc.Clientset.CoreV1().PersistentVolumeClaims(namespace).Get(name, metav1.GetOptions{})
}

// GetPVCs will get all PVCs from the k8s cluster
func (kc *KubeCtl) GetPVCs(namespace string) ([]*corev1.PersistentVolumeClaim, error) {
	pvcs := []*corev1.PersistentVolumeClaim{}
	pvcsList, err := kc.Clientset.CoreV1().PersistentVolumeClaims(namespace).List(metav1.ListOptions{})
	if err != nil {
		return pvcs, err
	}
	for _, p := range pvcsList.Items {
		pvcs = append(pvcs, &p)
	}
	return pvcs, nil
}

// CreatePVC will create the PVC by the PVC object
func (kc *KubeCtl) CreatePVC(pvc *corev1.PersistentVolumeClaim, namespace string) (*corev1.PersistentVolumeClaim, error) {
	return kc.Clientset.CoreV1().PersistentVolumeClaims(namespace).Create(pvc)
}

// DeletePVC will delete the PVC by the PVC name
func (kc *KubeCtl) DeletePVC(name string, namespace string) error {
	options := metav1.DeleteOptions{}
	return kc.Clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(name, &options)
}
