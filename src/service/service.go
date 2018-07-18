package service

import (
	"fmt"
	"regexp"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func checkNameValidation(name string) bool {
	re := regexp.MustCompile(`[a-z0-9]([-a-z0-9]*[a-z0-9])`)
	return re.MatchString(name)
}

func CheckServiceParameter(sp *serviceprovider.Container, service *entity.Service) error {
	session := sp.Mongo.NewSession()
	defer session.Close()

	//Check service name validation
	if !checkNameValidation(service.Name) {
		return fmt.Errorf("Service Name: %s is invalid value", service.Name)
	}

	//Check the service port name validation
	for _, port := range service.Ports {
		if !checkNameValidation(port.Name) {
			return fmt.Errorf("Port Name: %s is invalid value", port.Name)
		}
	}
	return nil
}

func CreateService(sp *serviceprovider.Container, service *entity.Service) error {
	session := sp.Mongo.NewSession()
	defer session.Close()

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

func DeleteService(sp *serviceprovider.Container, service *entity.Service) error {
	return sp.KubeCtl.DeleteService(service.Name, service.Namespace)
}
