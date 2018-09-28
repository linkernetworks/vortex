package configMap

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
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type ConfigMapTestSuite struct {
	suite.Suite
	sp *serviceprovider.Container
}

func (suite *ConfigMapTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	suite.sp = serviceprovider.NewForTesting(cf)
}

func (suite *ConfigMapTestSuite) TearDownSuite() {
}

func TestConfigMapSuite(t *testing.T) {
	suite.Run(t, new(ConfigMapTestSuite))
}

func (suite *ConfigMapTestSuite) TestCreateDeleteConfigMap() {
	configMapName := namesgenerator.GetRandomName(0)
	data := map[string]string{
		"firstData":  "awesome",
		"secondData": "cool",
	}
	configMap := &entity.ConfigMap{
		ID:        bson.NewObjectId(),
		Name:      configMapName,
		Namespace: "default",
		Data:      data,
	}

	err := CreateConfigMap(suite.sp, configMap)
	suite.NoError(err)

	err = DeleteConfigMap(suite.sp, configMap)
	suite.NoError(err)
}
