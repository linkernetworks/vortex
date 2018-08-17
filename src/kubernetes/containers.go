package kubernetes

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetContainer will get the container object by the container name
func (kc *KubeCtl) GetContainer(name string, podname string, namespace string) (*corev1.Container, error) {
	pod, _ := kc.Clientset.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
	for _, container := range pod.Spec.Containers {
		if container.Name == name {
			return &container, nil
		}
	}
	return nil, fmt.Errorf("can not find container")
}
