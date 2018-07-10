package kubernetes

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Get the service object by the service name
func (kc *KubeCtl) GetService(name string) (*corev1.Service, error) {
	return kc.Clientset.CoreV1().Services(kc.Namespace).Get(name, metav1.GetOptions{})
}

//Get all services from the k8s cluster
func (kc *KubeCtl) GetServices() ([]*corev1.Service, error) {
	services := []*corev1.Service{}
	servicesList, err := kc.Clientset.CoreV1().Services(kc.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return services, err
	}
	for _, s := range servicesList.Items {
		services = append(services, &s)
	}
	return services, nil
}

//Create the service by the service object
func (kc *KubeCtl) CreateService(service *corev1.Service) (*corev1.Service, error) {
	return kc.Clientset.CoreV1().Services(kc.Namespace).Create(service)
}

//Delete the service by the service name
func (kc *KubeCtl) DeleteService(name string) error {
	options := metav1.DeleteOptions{}
	return kc.Clientset.CoreV1().Services(kc.Namespace).Delete(name, &options)
}
