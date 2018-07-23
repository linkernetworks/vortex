package service

import (
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// CreateService will create service by serviceprovider container
func CreateService(sp *serviceprovider.Container, service *entity.Service) error {
	var serviceType corev1.ServiceType
	switch service.Type {
	default:
	case "ClusterIP":
		serviceType = corev1.ServiceTypeClusterIP
	case "NodePort":
		serviceType = corev1.ServiceTypeNodePort
	}

	var ports []corev1.ServicePort
	for _, port := range service.Ports {
		ports = append(ports, corev1.ServicePort{
			Name:       port.Name,
			Port:       port.Port,
			TargetPort: intstr.FromInt(port.TargetPort),
			NodePort:   port.NodePort,
		})
	}

	s := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: service.Name,
		},
		Spec: corev1.ServiceSpec{
			Type:     serviceType,
			Selector: service.Selector,
			Ports:    ports,
		},
	}
	_, err := sp.KubeCtl.CreateService(&s, service.Namespace)
	return err
}

// DeleteService willl delete service
func DeleteService(sp *serviceprovider.Container, service *entity.Service) error {
	return sp.KubeCtl.DeleteService(service.Name, service.Namespace)
}
