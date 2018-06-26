package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/redis"
	"github.com/linkernetworks/vortex/src/prometheus"
	"k8s.io/client-go/tools/clientcmd"
)

type Config struct {
	Redis      *redis.RedisConfig           `json:"redis"`
	Mongo      *mongo.MongoConfig           `json:"mongo"`
	Prometheus *prometheus.PrometheusConfig `json:"prometheus"`
	Logger     logger.LoggerConfig          `json:"logger"`

	// the version settings of the current application
	Version string `json:"version"`
}

func Read(path string) (c Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		return c, fmt.Errorf("Failed to open the config file: %v\n", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&c); err != nil {
		return c, fmt.Errorf("Failed to load the config file: %v\n", err)
	}

	// FIXME, we need to find a way to test the fakeclient evne if we don't install the k8s
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	c.Kubernetes, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return c, fmt.Errorf("Failed to open the kubernetes config file: %v\n", err)
	}

	return c, nil
}

func MustRead(path string) Config {
	c, err := Read(path)
	if err != nil {
		panic(err)
	}
	return c
}
