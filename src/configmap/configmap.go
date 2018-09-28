package configmap

import (
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateConfigMap will create configMap by serviceprovider container
func CreateConfigMap(sp *serviceprovider.Container, configMap *entity.ConfigMap) error {
	n := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMap.Name,
			Namespace: configMap.Namespace,
		},
		Data: configMap.Data,
	}
	_, err := sp.KubeCtl.CreateConfigMap(&n, configMap.Namespace)
	return err
}

// DeleteConfigMap will delete configMap
func DeleteConfigMap(sp *serviceprovider.Container, configMap *entity.ConfigMap) error {
	return sp.KubeCtl.DeleteConfigMap(configMap.Name, configMap.Namespace)
}
