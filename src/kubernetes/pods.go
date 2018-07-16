package kubernetes

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Get the pod object by the pod name
func (kc *KubeCtl) GetPod(name string, namespace string) (*corev1.Pod, error) {
	if namespace == "" {
		namespace = kc.Namespace
	}
	return kc.Clientset.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
}

//Get all pods from the k8s cluster
func (kc *KubeCtl) GetPods(namespace string) ([]*corev1.Pod, error) {
	if namespace == "" {
		namespace = kc.Namespace
	}
	pods := []*corev1.Pod{}
	podsList, err := kc.Clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return pods, err
	}
	for _, p := range podsList.Items {
		pods = append(pods, &p)
	}
	return pods, nil
}

//Create the pod by the pod object
func (kc *KubeCtl) CreatePod(pod *corev1.Pod, namespace string) (*corev1.Pod, error) {
	if namespace == "" {
		namespace = kc.Namespace
	}
	return kc.Clientset.CoreV1().Pods(namespace).Create(pod)
}

//Delete the pod by the pod name
func (kc *KubeCtl) DeletePod(name string, namespace string) error {
	if namespace == "" {
		namespace = kc.Namespace
	}
	options := metav1.DeleteOptions{}
	return kc.Clientset.CoreV1().Pods(namespace).Delete(name, &options)
}

func (kc *KubeCtl) IsPodCompleted(pod *corev1.Pod) bool {
	switch pod.Status.Phase {
	case corev1.PodRunning, corev1.PodPending:
		return false
	default:
		return true
	}
	return true
}
