package kubernetes

import (
	"math/rand"
	"testing"
	"time"

	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"
	"k8s.io/api/autoscaling/v2beta1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

type KubeCtlAutoscalerTestSuite struct {
	suite.Suite
	kubectl    *KubeCtl
	fakeclient *fakeclientset.Clientset
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (suite *KubeCtlAutoscalerTestSuite) SetupSuite() {
	suite.fakeclient = fakeclientset.NewSimpleClientset()
	suite.kubectl = New(suite.fakeclient)
}

func (suite *KubeCtlAutoscalerTestSuite) TearDownSuite() {}

func (suite *KubeCtlAutoscalerTestSuite) TestCreateAutoscaler() {
	namespace := "default"
	var replicas, minReplicas, targetAverageUtilization int32
	replicas = 3
	minReplicas = 3
	targetAverageUtilization = 30
	deploymentName := namesgenerator.GetRandomName(0)

	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
		},
		Status: appsv1.DeploymentStatus{},
	}

	autoscaler := v2beta1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			// use deployment name to name autoscaler's name
			Name:      deploymentName,
			Namespace: namespace,
		},
		Spec: v2beta1.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: v2beta1.CrossVersionObjectReference{
				Kind: "Deployment",
				Name: deploymentName,
			},
			MinReplicas: &minReplicas,
			MaxReplicas: 3,
			Metrics: []v2beta1.MetricSpec{
				{
					Type: v2beta1.ResourceMetricSourceType,
					Resource: &v2beta1.ResourceMetricSource{
						Name:                     corev1.ResourceCPU,
						TargetAverageUtilization: &targetAverageUtilization,
					},
				},
			},
		},
	}

	ret, err := suite.kubectl.CreateDeployment(&deployment, namespace)
	suite.NoError(err)
	suite.NotNil(ret)

	deploymentResponse, err := suite.kubectl.GetDeployment(deploymentName, namespace)
	suite.NoError(err)
	suite.NotNil(deploymentResponse)

	autoscalerResponse, err := suite.kubectl.CreateAutoscaler(&autoscaler, namespace)
	suite.NoError(err)
	suite.NotNil(autoscalerResponse)

	defer suite.kubectl.DeleteAutoscaler(deploymentName, namespace)
	defer suite.kubectl.DeleteDeployment(deploymentName, namespace)
}

func (suite *KubeCtlAutoscalerTestSuite) TestGetAutoscaler() {
	namespace := "default"
	var replicas, minReplicas, targetAverageUtilization int32
	replicas = 3
	minReplicas = 3
	targetAverageUtilization = 30
	deploymentName := namesgenerator.GetRandomName(0)

	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
		},
		Status: appsv1.DeploymentStatus{},
	}

	autoscaler := v2beta1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			// use deployment name to name autoscaler's name
			Name:      deploymentName,
			Namespace: namespace,
		},
		Spec: v2beta1.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: v2beta1.CrossVersionObjectReference{
				Kind: "Deployment",
				Name: deploymentName,
			},
			MinReplicas: &minReplicas,
			MaxReplicas: 3,
			Metrics: []v2beta1.MetricSpec{
				{
					Type: v2beta1.ResourceMetricSourceType,
					Resource: &v2beta1.ResourceMetricSource{
						Name:                     corev1.ResourceCPU,
						TargetAverageUtilization: &targetAverageUtilization,
					},
				},
			},
		},
	}

	ret, err := suite.kubectl.CreateDeployment(&deployment, namespace)
	suite.NoError(err)
	suite.NotNil(ret)

	deploymentResponse, err := suite.kubectl.GetDeployment(deploymentName, namespace)
	suite.NoError(err)
	suite.NotNil(deploymentResponse)

	autoscalerResponse, err := suite.kubectl.CreateAutoscaler(&autoscaler, namespace)
	suite.NoError(err)
	suite.NotNil(autoscalerResponse)

	autoscalerGetResponse, err := suite.kubectl.GetAutoscaler(autoscaler.Name, namespace)
	suite.NoError(err)
	suite.NotNil(autoscalerGetResponse)

	defer suite.kubectl.DeleteAutoscaler(deploymentName, namespace)
	defer suite.kubectl.DeleteDeployment(deploymentName, namespace)
}

func (suite *KubeCtlAutoscalerTestSuite) TestDeleteAutoscaler() {
	namespace := "default"
	var replicas, minReplicas, targetAverageUtilization int32
	replicas = 3
	minReplicas = 3
	targetAverageUtilization = 30
	deploymentName := namesgenerator.GetRandomName(0)

	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
		},
		Status: appsv1.DeploymentStatus{},
	}

	autoscaler := v2beta1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			// use deployment name to name autoscaler's name
			Name:      deploymentName,
			Namespace: namespace,
		},
		Spec: v2beta1.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: v2beta1.CrossVersionObjectReference{
				Kind: "Deployment",
				Name: deploymentName,
			},
			MinReplicas: &minReplicas,
			MaxReplicas: 3,
			Metrics: []v2beta1.MetricSpec{
				{
					Type: v2beta1.ResourceMetricSourceType,
					Resource: &v2beta1.ResourceMetricSource{
						Name:                     corev1.ResourceCPU,
						TargetAverageUtilization: &targetAverageUtilization,
					},
				},
			},
		},
	}
	ret, err := suite.kubectl.CreateDeployment(&deployment, namespace)
	suite.NoError(err)
	suite.NotNil(ret)

	deploy, err := suite.kubectl.GetDeployment(deploymentName, namespace)
	suite.NoError(err)
	suite.NotNil(deploy)
	suite.Equal(replicas, *deploy.Spec.Replicas)

	autoscalerResponse, err := suite.kubectl.CreateAutoscaler(&autoscaler, namespace)
	suite.NoError(err)
	suite.NotNil(autoscalerResponse)

	err = suite.kubectl.DeleteDeployment(deploymentName, namespace)
	suite.NoError(err)
	deploy, err = suite.kubectl.GetDeployment(deploymentName, namespace)
	suite.Error(err)
	suite.Nil(deploy)
	defer suite.kubectl.DeleteAutoscaler(deploymentName, namespace)
}

func TestAutoscalerTestSuite(t *testing.T) {
	suite.Run(t, new(KubeCtlAutoscalerTestSuite))
}
