package serviceprovider

import (
	"github.com/linkernetworks/config"
	"github.com/linkernetworks/logger"

	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/redis"
)

type Container struct {
	Config config.Config
	Redis  *redis.Service
	Mongo  *mongo.Service
}

type ServiceDiscoverResponse struct {
	Container map[string]Service `json:"services"`
}

type Service interface{}

func New(cf config.Config) *Container {
	// setup logger configuration
	logger.Setup(cf.Logger)

	logger.Infof("Connecting to redis: %s", cf.Redis.Addr())
	redisService := redis.New(cf.Redis)

	logger.Infof("Connecting to mongodb: %s", cf.Mongo.Url)
	mongo := mongo.New(cf.Mongo.Url)

	sp := &Container{
		Config: cf,
		Redis:  redisService,
		Mongo:  mongo,
	}

	return sp
}

func NewContainer(configPath string) *Container {
	cf := config.MustRead(configPath)
	return New(cf)
}
