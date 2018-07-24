package kubernetes

import (
	"k8s.io/client-go/kubernetes"
)


// KubeCtl object is used to interact with the kubernetes cluster.
// Use the export function New to Get a KubeCtl object.
type KubeCtl struct {
	Clientset kubernetes.Interface
}


// New is the API to New a kubectl object and you need to pass two parameters
// 1. The kubernetes clientset object from the client-go library. You can also use the fake-client for testing
// 2. The namespace of the kubernetes you want to manipulate
func New(clientset kubernetes.Interface) *KubeCtl {
	return &KubeCtl{
		Clientset: clientset,
	}
}
