package kubernetes

import (
	v1beta2 "k8s.io/api/autoscaling/v2beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateAutoscaler will create a autoscaler
func (kc *KubeCtl) CreateAutoscaler(autoscaler *v1beta2.HorizontalPodAutoscaler, namespace string) (*v1beta2.HorizontalPodAutoscaler, error) {
	return kc.Clientset.AutoscalingV2beta1().HorizontalPodAutoscalers(namespace).Create(autoscaler)
}

// GetAutoscaler will get a autoscaler
func (kc *KubeCtl) GetAutoscaler(name string, namespace string) (*v1beta2.HorizontalPodAutoscaler, error) {
	return kc.Clientset.AutoscalingV2beta1().HorizontalPodAutoscalers(namespace).Get(name, metav1.GetOptions{})
}

// DeleteAutoscaler will delete a autoscaler
func (kc *KubeCtl) DeleteAutoscaler(name string, namespace string) error {
	propagation := metav1.DeletePropagationForeground
	return kc.Clientset.AutoscalingV2beta1().HorizontalPodAutoscalers(namespace).Delete(name, &metav1.DeleteOptions{PropagationPolicy: &propagation})
}
