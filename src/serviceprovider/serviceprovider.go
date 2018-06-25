package serviceprovider

import (
	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/vortex/src/config"

	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/redis"
	kubeCtl "github.com/linkernetworks/vortex/src/kubernetes"
	"k8s.io/client-go/kubernetes"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

type Container struct {
	Config  config.Config
	Redis   *redis.Service
	Mongo   *mongo.Service
	KubeCtl *kubeCtl.KubeCtl
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

	var clientset kubernetes.Interface

	if cf.Kubernetes != nil {
		clientset = kubernetes.NewForConfigOrDie(cf.Kubernetes)
	} else {
		clientset = fakeclientset.NewSimpleClientset()
	}

	sp := &Container{
		Config:  cf,
		Redis:   redisService,
		Mongo:   mongo,
		KubeCtl: kubeCtl.New(clientset, "default"),
	}

	return sp
}

func NewContainer(configPath string) *Container {
	cf := config.MustRead(configPath)
	return New(cf)
}
