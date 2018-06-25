package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/redis"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Config struct {
	Redis  *redis.RedisConfig  `json:"redis"`
	Mongo  *mongo.MongoConfig  `json:"mongo"`
	Logger logger.LoggerConfig `json:"logger"`

	// the version settings of the current application
	Version string `json:"version"`

	Kubernetes *rest.Config `json:"kubernetes"`
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
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	c.Kubernetes, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		c.Kubernetes, err = rest.InClusterConfig()
		if err != nil {
			return c, fmt.Errorf("Failed to open the kubernetes config file: %v\n", err)
		}
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
