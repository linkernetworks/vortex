package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/vortex/src/prometheusprovider"
	"github.com/linkernetworks/vortex/src/registry"
)

// Config is the structure for vortex
type Config struct {
	Mongo      *mongo.MongoConfig                   `json:"mongo"`
	Prometheus *prometheusprovider.PrometheusConfig `json:"prometheus"`
	Registry   *registry.Config                     `json:"registry"`
	Kubernetes *kubernetesConfig                    `json:"kubernetes"`
	Logger     logger.LoggerConfig                  `json:"logger"`

	// the version settings of the current application
	Version string `json:"version"`
}

// Read will read config file
func Read(path string) (c Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		return c, fmt.Errorf("Failed to open the config file: %v", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&c); err != nil {
		return c, fmt.Errorf("Failed to load the config file: %v", err)
	}

	return c, nil
}

// MustRead will read config path
func MustRead(path string) Config {
	c, err := Read(path)
	if err != nil {
		panic(err)
	}
	return c
}
