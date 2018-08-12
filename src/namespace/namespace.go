package namespace

import (
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateNamespace will create namespace by serviceprovider container
func CreateNamespace(sp *serviceprovider.Container, namespace *entity.Namespace) error {
	n := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   namespace.Name,
		},
	}
	_, err := sp.KubeCtl.CreateNamespace(&n)
	return err
}

// DeleteNamespace willl delete namespace
func DeleteNamespace(sp *serviceprovider.Container, namespace *entity.Namespace) error {
	return sp.KubeCtl.DeleteNamespace(namespace.Name)
}
