package service

import (
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type ServiceTestSuite struct {
	suite.Suite
	sp *serviceprovider.Container
}

func (suite *ServiceTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	suite.sp = serviceprovider.NewForTesting(cf)
}

func (suite *ServiceTestSuite) TearDownSuite() {
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) TestCheckServiceParameter() {
	service := &entity.Service{
		ID:   bson.NewObjectId(),
		Name: namesgenerator.GetRandomName(0),
	}

	err := CheckServiceParameter(suite.sp, service)
	suite.NoError(err)
}

func (suite *ServiceTestSuite) TestCheckServiceParameterFail() {
	testCases := []struct {
		caseName string
		service  *entity.Service
	}{
		{
			"InvalidServiceName", &entity.Service{
				ID:   bson.NewObjectId(),
				Name: "~!@#$%^&*()",
			},
		},
		{
			"InvalidServicePortName", &entity.Service{
				ID:   bson.NewObjectId(),
				Name: namesgenerator.GetRandomName(0),
				Ports: []entity.ServicePort{
					{
						Name: "~!@#$%^&*()",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.caseName, func(t *testing.T) {
			err := CheckServiceParameter(suite.sp, tc.service)
			suite.Error(err)
		})
	}
}

func (suite *ServiceTestSuite) TestCreateService() {
	selector := map[string]string{
		"podname": "awesome",
	}
	ports := []entity.ServicePort{
		{
			Name:       namesgenerator.GetRandomName(0),
			Port:       int32(80),
			TargetPort: 80,
			NodePort:   int32(30000),
		},
	}

	serviceName := namesgenerator.GetRandomName(0)
	service := &entity.Service{
		ID:        bson.NewObjectId(),
		Name:      serviceName,
		Namespace: "default",
		Type:      "NodePort",
		Selector:  selector,
		Ports:     ports,
	}

	err := CreateService(suite.sp, service)
	suite.NoError(err)

	err = DeleteService(suite.sp, service)
	suite.NoError(err)
}
