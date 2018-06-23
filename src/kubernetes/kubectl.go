package kubernetes

import (
	"k8s.io/client-go/kubernetes"
)

type KubeCtl struct {
	Clientset kubernetes.Interface
	Namespace string
}

func New(clientset kubernetes.Interface, namespace string) *KubeCtl {
	return &KubeCtl{
		Clientset: clientset,
		Namespace: namespace,
	}
}
