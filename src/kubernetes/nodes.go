package kubernetes

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Get the node object by the node name
func (kc *KubeCtl) GetNode(name string) (*corev1.Node, error) {
	return kc.Clientset.CoreV1().Nodes().Get(name, metav1.GetOptions{})
}

//Get all nodes from the k8s cluster
func (kc *KubeCtl) GetNodes() ([]*corev1.Node, error) {
	nodes := []*corev1.Node{}
	nodesList, err := kc.Clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nodes, err
	}
	for _, n := range nodesList.Items {
		nodes = append(nodes, &n)
	}
	return nodes, nil
}

//Get the external IP address of node
func (kc *KubeCtl) GetNodeExternalIP(name string) (string, error) {
	node, err := kc.GetNode(name)
	if err != nil {
		return "", err
	}
	var nodeIP string
	for _, addr := range node.Status.Addresses {
		if addr.Type == "ExternalIP" {
			nodeIP = addr.Address
			break
		}
	}
	return nodeIP, nil
}

//Get the internal IP address of node
func (kc *KubeCtl) GetNodeInternalIP(name string) (string, error) {
	node, err := kc.GetNode(name)
	if err != nil {
		return "", err
	}
	var nodeIP string
	for _, addr := range node.Status.Addresses {
		if addr.Type == "InternalIP" {
			nodeIP = addr.Address
			break
		}
	}
	return nodeIP, nil
}
