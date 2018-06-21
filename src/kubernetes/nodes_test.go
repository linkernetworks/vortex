package kubernetes

import (
	"github.com/linkernetworks/kubeconfig"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes"
	"testing"
)

func TestGetNodeFail(t *testing.T) {
	config, err := kubeconfig.Load("")
	assert.NoError(t, err)

	clientset, err := kubernetes.NewForConfig(config)
	assert.NoError(t, err)

	_, err = GetNode(clientset, "UnKnown_Name")
	assert.Error(t, err)
}

func TestGetNodes(t *testing.T) {
	config, err := kubeconfig.Load("")
	assert.NoError(t, err)

	clientset, err := kubernetes.NewForConfig(config)
	assert.NoError(t, err)

	_, err = GetNodes(clientset)
	assert.NoError(t, err)
}
