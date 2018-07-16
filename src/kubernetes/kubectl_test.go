package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"

	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

func TestNewKubeCtl(t *testing.T) {
	clientset := fakeclientset.NewSimpleClientset()
	kubectl := New(clientset)
	assert.NotNil(t, kubectl)
}
