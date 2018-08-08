package kubernetes

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetNamespace will get the namespace object by the namespace name
func (kc *KubeCtl) GetNamespace(name string) (*corev1.Namespace, error) {
	return kc.Clientset.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
}

// GetNamespaces will get all namespaces from the k8s cluster
func (kc *KubeCtl) GetNamespaces() ([]*corev1.Namespace, error) {
	namespaces := []*corev1.Namespace{}
	namespacesList, err := kc.Clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return namespaces, err
	}
	for _, n := range namespacesList.Items {
		namespaces = append(namespaces, &n)
	}
	return namespaces, nil
}

// CreateNamespace will create the namespace by the namespace object
func (kc *KubeCtl) CreateNamespace(namespace *corev1.Namespace) (*corev1.Namespace, error) {
	return kc.Clientset.CoreV1().Namespaces().Create(namespace)
}

// DeleteNamespace will delete the namespace by the namespace name
func (kc *KubeCtl) DeleteNamespace(name string) error {
	deletePolicy := metav1.DeletePropagationBackground
	return kc.Clientset.CoreV1().Namespaces().Delete(name, &metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
}
