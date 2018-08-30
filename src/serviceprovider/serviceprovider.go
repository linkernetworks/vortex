package serviceprovider

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/prometheusprovider"

	"github.com/linkernetworks/mongo"
	kubeCtl "github.com/linkernetworks/vortex/src/kubernetes"

	"gopkg.in/go-playground/validator.v9"
	"k8s.io/client-go/kubernetes"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Container is the structure for container
type Container struct {
	Config     config.Config
	Mongo      *mongo.Service
	Prometheus *prometheusprovider.Service
	KubeCtl    *kubeCtl.KubeCtl
	Validator  *validator.Validate
}

// ServiceDiscoverResponse is the structure for Service Discover Response
type ServiceDiscoverResponse struct {
	Container map[string]Service `json:"services"`
}

// Service is the interface
type Service interface{}

// New will create container
func New(cf config.Config) *Container {
	// setup logger configuration
	logger.Setup(cf.Logger)

	logger.Infof("Connecting to mongodb: %s", cf.Mongo.Url)
	mongo := mongo.New(cf.Mongo.Url)

	logger.Infof("Connecting to prometheus: %s", cf.Prometheus.URL)
	prometheus := prometheusprovider.New(cf.Prometheus.URL)

	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	k8s, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		k8s, err = rest.InClusterConfig()
		if err != nil {
			panic(fmt.Errorf("Load the kubernetes config fail"))
		}
	}

	clientset := kubernetes.NewForConfigOrDie(k8s)

	validate := validator.New()
	// Register validation for kubernetes name
	validate.RegisterValidation("k8sname", checkNameValidation)

	sp := &Container{
		Config:     cf,
		Mongo:      mongo,
		Prometheus: prometheus,
		KubeCtl:    kubeCtl.New(clientset),
		Validator:  validate,
	}

	if err := createDefaultUser(sp.Mongo); err != nil {
		logger.Infof("Create Default admin user failed: %v", err)
	}
	return sp
}

// NewForTesting will test container for creating a container
func NewForTesting(cf config.Config) *Container {
	// setup logger configuration
	logger.Setup(cf.Logger)

	logger.Infof("Connecting to mongodb: %s", cf.Mongo.Url)
	mongo := mongo.New(cf.Mongo.Url)

	logger.Infof("Connecting to prometheus: %s", cf.Prometheus.URL)
	prometheus := prometheusprovider.New(cf.Prometheus.URL)

	clientset := fakeclientset.NewSimpleClientset()

	validate := validator.New()
	// Register validation for kubernetes name
	validate.RegisterValidation("k8sname", checkNameValidation)

	sp := &Container{
		Config:     cf,
		Mongo:      mongo,
		Prometheus: prometheus,
		KubeCtl:    kubeCtl.New(clientset),
		Validator:  validate,
	}

	return sp
}

// NewContainer will new a container
func NewContainer(configPath string) *Container {
	cf := config.MustRead(configPath)
	return New(cf)
}
