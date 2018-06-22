package serviceprovider

import (
	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/vortex/src/config"

	"github.com/linkernetworks/kubeconfig"
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/redis"
	"k8s.io/client-go/kubernetes"
)

type Container struct {
	Config     config.Config
	Redis      *redis.Service
	Mongo      *mongo.Service
	Kubernetes kubernetes.Interface
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

	k8sConfig, err := kubeconfig.Load(cf.Kubernetes)
	if err != nil {
		logger.Warnln("kubernetes service is not loaded: kubernetes config is not defined.")
	} else {
		clientset, err := kubernetes.NewForConfig(k8sConfig)
		if err != nil {
			logger.Warnln("kubernetes service is not loaded: clientset creates fail.")
		}
		sp.Kubernetes = clientset
	}

	return sp
}

func NewContainer(configPath string) *Container {
	cf := config.MustRead(configPath)
	return New(cf)
}
