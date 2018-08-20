package kubernetes

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateDeployment will get the external IP address of node
func (kc *KubeCtl) CreateDeployment(deployment *appsv1.Deployment, namespace string) (*appsv1.Deployment, error) {
	return kc.Clientset.AppsV1().Deployments(namespace).Create(deployment)
}

// GetDeployment will get deploy
func (kc *KubeCtl) GetDeployment(name string, namespace string) (*appsv1.Deployment, error) {
	return kc.Clientset.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
}

// GetDeployments will get deploys
func (kc *KubeCtl) GetDeployments(namespace string) ([]*appsv1.Deployment, error) {
	deployments := []*appsv1.Deployment{}
	deploymentsList, err := kc.Clientset.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
	if err != nil {
		return deployments, err
	}
	for _, d := range deploymentsList.Items {
		deployments = append(deployments, &d)
	}
	return deployments, nil
}

// DeleteDeployment will delete deploy
func (kc *KubeCtl) DeleteDeployment(name string, namespace string) error {
	propagation := metav1.DeletePropagationForeground
	return kc.Clientset.AppsV1().Deployments(namespace).Delete(name, &metav1.DeleteOptions{PropagationPolicy: &propagation})
}
