package kubernetes

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetNode(clientset *kubernetes.Clientset, name string) (*v1.Node, error) {
	return clientset.CoreV1().Nodes().Get(name, metav1.GetOptions{})
}

func GetNodes(clientset *kubernetes.Clientset) ([]*v1.Node, error) {
	nodes := []*v1.Node{}
	nodesList, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nodes, err
	}

	for _, n := range nodesList.Items {
		nodes = append(nodes, &n)
	}

	return nodes, nil
}
