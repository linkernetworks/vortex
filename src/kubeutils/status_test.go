package kubeutils

import (
	"math/rand"
	"testing"
	"time"

	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/moby/moby/pkg/namesgenerator"

	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2/bson"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type StatusTestSuite struct {
	suite.Suite
	sp *serviceprovider.Container
}

func (suite *StatusTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	suite.sp = serviceprovider.NewForTesting(cf)
}

func (suite *StatusTestSuite) TearDownSuite() {
}

func TestStatusSuite(t *testing.T) {
	suite.Run(t, new(StatusTestSuite))
}

func (suite *StatusTestSuite) TestGetNonCompletedPods() {
	namespace := "default"
	session := suite.sp.Mongo.NewSession()
	defer session.Close()
	networkName := namesgenerator.GetRandomName(0)

	pods := []entity.Pod{
		{
			ID:   bson.NewObjectId(),
			Name: namesgenerator.GetRandomName(0),
			Volumes: []entity.PodVolume{
				{},
			},
			Networks: []entity.PodNetwork{
				{Name: networkName},
			},
		},
		{
			ID:   bson.NewObjectId(),
			Name: namesgenerator.GetRandomName(0),
			Volumes: []entity.PodVolume{
				{},
			},
			Networks: []entity.PodNetwork{
				{Name: networkName},
			},
		},
	}

	for _, pod := range pods {
		session.Insert(entity.PodCollectionName, pod)
		defer session.Remove(entity.PodCollectionName, "name", pod.Name)
	}

	//Create the pod via kubectl
	suite.sp.KubeCtl.CreatePod(&corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: pods[0].Name,
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
	}, namespace)
	suite.sp.KubeCtl.CreatePod(&corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: pods[1].Name,
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
	}, namespace)

	ret, err := GetNonCompletedPods(suite.sp, bson.M{"networks.name": networkName})
	suite.Equal(len(pods), len(ret))
	suite.NoError(err)
	ret, err = GetNonCompletedPods(suite.sp, bson.M{"networks.name": namesgenerator.GetRandomName(1)})
	suite.Equal(0, len(ret))
	suite.NoError(err)
}
