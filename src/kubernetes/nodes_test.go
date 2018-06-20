package kubernetes

import (
	"github.com/linkernetworks/config"
	"github.com/linkernetworks/service/kubernetes"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetNodeFail(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_K8S"); !ok {
		t.SkipNow()
	}

	kubernetes := kubernetes.NewFromConfig(&config.KubernetesConfig{
		Namespace: "default",
	})
	clientset, err := kubernetes.NewClientset()
	assert.NoError(t, err)

	_, err = GetNode(clientset, "UnKnown_Name")
	assert.Error(t, err)
}

func TestGetNodes(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_K8S"); !ok {
		t.SkipNow()
	}

	kubernetes := kubernetes.NewFromConfig(&config.KubernetesConfig{
		Namespace: "default",
	})
	clientset, err := kubernetes.NewClientset()
	assert.NoError(t, err)

	nodes, err := GetNodes(clientset)
	assert.NoError(t, err)
}
