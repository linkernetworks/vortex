package kubernetes

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Get the pod object by the pod name
func (kc *KubeCtl) GetPod(name string) (*corev1.Pod, error) {
	return kc.Clientset.CoreV1().Pods(kc.Namespace).Get(name, metav1.GetOptions{})
}

//Get all pods from the k8s cluster
func (kc *KubeCtl) GetPods() ([]*corev1.Pod, error) {
	pods := []*corev1.Pod{}
	podsList, err := kc.Clientset.CoreV1().Pods(kc.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return pods, err
	}
	for _, p := range podsList.Items {
		pods = append(pods, &p)
	}
	return pods, nil
}

//Create the pod by the pod object
func (kc *KubeCtl) CreatePod(pod *corev1.Pod) (*corev1.Pod, error) {
	return kc.Clientset.CoreV1().Pods(kc.Namespace).Create(pod)
}

//Delete the pod by the pod name
func (kc *KubeCtl) DeletePod(name string) error {
	options := metav1.DeleteOptions{}
	return kc.Clientset.CoreV1().Pods(kc.Namespace).Delete(name, &options)
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
