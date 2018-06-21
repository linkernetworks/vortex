package serviceprovider

import (
	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/vortex/src/config"

	"github.com/linkernetworks/kubeconfig"
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/redis"
	"k8s.io/client-go/rest"
)

type Container struct {
	Config     config.Config
	Redis      *redis.Service
	Mongo      *mongo.Service
	Kubernetes *rest.Config
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

	k8sConfig, _ := kubeconfig.Load(cf.Kubernetes)

	sp := &Container{
		Config:     cf,
		Redis:      redisService,
		Mongo:      mongo,
		Kubernetes: k8sConfig,
	}

	return sp
}

func NewContainer(configPath string) *Container {
	cf := config.MustRead(configPath)
	return New(cf)
}
