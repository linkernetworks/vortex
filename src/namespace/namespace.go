package namespace

import (
	"fmt"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// CreateNamespace will create namespace by serviceprovider container
func CreateNamespace(sp *serviceprovider.Container, namespace *entity.Namespace) error {
	n := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace.Name,
		},
	}
	_, err := sp.KubeCtl.CreateNamespace(&n)
	return err
}

// DeleteNamespace will delete namespace
func DeleteNamespace(sp *serviceprovider.Container, namespace *entity.Namespace) error {
	deploys, _ := sp.KubeCtl.GetDeployments(namespace.Name)
	svcs, _ := sp.KubeCtl.GetServices(namespace.Name)
	pvcs, _ := sp.KubeCtl.GetPVCs(namespace.Name)

	if len(deploys) != 0 || len(svcs) != 0 || len(pvcs) != 0 {
		return errors.NewForbidden(schema.GroupResource{Group: "none", Resource: "Namespace"}, namespace.Name, fmt.Errorf("Still have some resource under namespace %v", namespace.Name))
	}

	return sp.KubeCtl.DeleteNamespace(namespace.Name)
}
