package kubernetes

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetConfigMap will get the configMap object by the configMap name
func (kc *KubeCtl) GetConfigMap(name string, namespace string) (*corev1.ConfigMap, error) {
	return kc.Clientset.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
}

// GetConfigMaps will get all configMaps from the k8s cluster
func (kc *KubeCtl) GetConfigMaps(namespace string) ([]*corev1.ConfigMap, error) {
	configMaps := []*corev1.ConfigMap{}
	configMapsList, err := kc.Clientset.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{})
	if err != nil {
		return configMaps, err
	}
	for _, n := range configMapsList.Items {
		configMaps = append(configMaps, &n)
	}
	return configMaps, nil
}

// CreateConfigMap will create the configMap by the configMap object
func (kc *KubeCtl) CreateConfigMap(configMap *corev1.ConfigMap, namespace string) (*corev1.ConfigMap, error) {
	return kc.Clientset.CoreV1().ConfigMaps(namespace).Create(configMap)
}

// DeleteConfigMap will delete the configMap by the configMap name
func (kc *KubeCtl) DeleteConfigMap(name string, namespace string) error {
	foreground := metav1.DeletePropagationForeground
	return kc.Clientset.CoreV1().ConfigMaps(namespace).Delete(name, &metav1.DeleteOptions{PropagationPolicy: &foreground})
}
