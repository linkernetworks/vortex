package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/redis"
)

type Config struct {
	Redis  *redis.RedisConfig  `json:"redis"`
	Mongo  *mongo.MongoConfig  `json:"mongo"`
	Logger logger.LoggerConfig `json:"logger"`

	// the version settings of the current application
	Version string `json:"version"`

	Kubernetes string `json:"kubernetes"`
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
	return c, nil
}

func MustRead(path string) Config {
	c, err := Read(path)
	if err != nil {
		panic(err)
	}
	return c
}
