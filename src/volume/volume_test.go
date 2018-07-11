package volume

import (
	"math/rand"
	"testing"
	"time"

	//"github.com/linkernetworks/mongo"
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

type VolumeTestSuite struct {
	suite.Suite
	sp *serviceprovider.Container
}

func (suite *VolumeTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	suite.sp = serviceprovider.NewForTesting(cf)
}

func (suite *VolumeTestSuite) TearDownSuite() {
}

func TestVolumeSuite(t *testing.T) {
	suite.Run(t, new(VolumeTestSuite))
}

func (suite *VolumeTestSuite) TestGetPVCInstance() {
	volume := &entity.Volume{
		ID:          bson.NewObjectId(),
		Name:        namesgenerator.GetRandomName(0),
		StorageName: namesgenerator.GetRandomName(0),
	}

	pvc := getPVCInstance(volume, namesgenerator.GetRandomName(0), namesgenerator.GetRandomName(0))
	suite.NotNil(pvc)
}

func (suite *VolumeTestSuite) TestGetStorageClassName() {
	session := suite.sp.Mongo.NewSession()
	defer session.Close()

	storage := entity.Storage{
		ID:               bson.NewObjectId(),
		Name:             namesgenerator.GetRandomName(0),
		StorageClassName: namesgenerator.GetRandomName(0),
	}

	session.Insert(entity.StorageCollectionName, &storage)
	defer session.Remove(entity.StorageCollectionName, "name", storage.Name)
	name, err := getStorageClassName(session, storage.Name)
	suite.NoError(err)
	suite.Equal(name, storage.StorageClassName)
}

func (suite *VolumeTestSuite) TestCreateVolume() {
	session := suite.sp.Mongo.NewSession()
	defer session.Close()
	storage := entity.Storage{
		ID:               bson.NewObjectId(),
		Name:             namesgenerator.GetRandomName(0),
		StorageClassName: namesgenerator.GetRandomName(0),
	}

	session.Insert(entity.StorageCollectionName, &storage)
	defer session.Remove(entity.StorageCollectionName, "name", storage.Name)

	volume := &entity.Volume{
		ID:          bson.NewObjectId(),
		Name:        namesgenerator.GetRandomName(0),
		StorageName: storage.Name,
	}

	err := CreateVolume(suite.sp, volume)
	suite.NoError(err)

	name := volume.GetPVCName()
	v, err := suite.sp.KubeCtl.GetPVC(name)
	suite.NoError(err)
	suite.NotNil(v)

	err = DeleteVolume(suite.sp, volume)
	suite.NoError(err)

	v, err = suite.sp.KubeCtl.GetPVC(name)
	suite.Error(err)
	suite.Nil(v)
}

func (suite *VolumeTestSuite) TestCreateVolumeFail() {
	volume := &entity.Volume{
		ID:          bson.NewObjectId(),
		Name:        namesgenerator.GetRandomName(0),
		StorageName: namesgenerator.GetRandomName(0),
	}

	err := CreateVolume(suite.sp, volume)
	suite.Error(err)
}

func (suite *VolumeTestSuite) TestDeleteVolumeFail() {
	volume := &entity.Volume{
		ID:   bson.NewObjectId(),
		Name: namesgenerator.GetRandomName(0),
	}

	session := suite.sp.Mongo.NewSession()
	defer session.Close()

	pods := []entity.Pod{
		{
			ID:   bson.NewObjectId(),
			Name: namesgenerator.GetRandomName(0),
			Volumes: []entity.PodVolume{
				{
					Name: volume.Name,
				},
			},
		},
		{
			ID:   bson.NewObjectId(),
			Name: namesgenerator.GetRandomName(0),
			Volumes: []entity.PodVolume{
				{
					Name: volume.Name,
				},
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
	})
	suite.sp.KubeCtl.CreatePod(&corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: pods[1].Name,
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
	})
	suite.sp.KubeCtl.CreatePod(&corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: namesgenerator.GetRandomName(0),
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
	})

	err := DeleteVolume(suite.sp, volume)
	suite.Error(err)
}
