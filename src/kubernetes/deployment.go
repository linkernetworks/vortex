package kubernetes

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Get the external IP address of node
func (kc *KubeCtl) CreateDeployment(deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return kc.Clientset.AppsV1().Deployments(kc.Namespace).Create(deployment)
}

func (kc *KubeCtl) GetDeployment(name string) (*appsv1.Deployment, error) {
	return kc.Clientset.AppsV1().Deployments(kc.Namespace).Get(name, metav1.GetOptions{})
}

func (kc *KubeCtl) GetDeployments() ([]*appsv1.Deployment, error) {
	deployments := []*appsv1.Deployment{}
	deploymentsList, err := kc.Clientset.AppsV1().Deployments(kc.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return deployments, err
	}
	for _, d := range deploymentsList.Items {
		deployments = append(deployments, &d)
	}
	return deployments, nil
}

func (kc *KubeCtl) DeleteDeployment(name string) error {
	return kc.Clientset.AppsV1().Deployments(kc.Namespace).Delete(name, &metav1.DeleteOptions{})
}
