package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"

	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

func TestNewKubeCtl(t *testing.T) {
	clientset := fakeclientset.NewSimpleClientset()
	namespace := "default"
	kubectl := New(clientset, namespace)
	assert.Equal(t, namespace, kubectl.Namespace)
	assert.NotNil(t, kubectl)
}

func TestChangeKubeCtlNamespace(t *testing.T) {
	clientset := fakeclientset.NewSimpleClientset()
	namespace := "default"
	kubectl := New(clientset, namespace)
	assert.Equal(t, namespace, kubectl.Namespace)
	assert.NotNil(t, kubectl)

	kubectl.ChangeNamespace("test")
}
