package serviceprovider

import (
	"github.com/linkernetworks/config"
	"github.com/linkernetworks/logger"

	"github.com/linkernetworks/gearman"
	"github.com/linkernetworks/influxdb"
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/redis"
)

type Container struct {
	Config   config.Config
	Redis    *redis.Service
	Mongo    *mongo.Service
	Gearman  *gearman.Service
	Influxdb *influxdb.InfluxdbService
}

type ServiceDiscoverResponse struct {
	Container map[string]Service `json:"services"`
}

type Service interface{}

func NewInfluxdbService(cf *influxdb.InfluxdbConfig) *influxdb.InfluxdbService {
	logger.Infof("Connecting to influxdb: %s", cf.Url)
	return &influxdb.InfluxdbService{Url: cf.Url}
}

func New(cf config.Config) *Container {
	// setup logger configuration
	logger.Setup(cf.Logger)

	logger.Infof("Connecting to redis: %s", cf.Redis.Addr())
	redisService := redis.New(cf.Redis)

	logger.Infof("Connecting to mongodb: %s", cf.Mongo.Url)
	mongo := mongo.New(cf.Mongo.Url)

	logger.Infof("Connecting to influxdb: %s", cf.Influxdb.Url)
	influx := &influxdb.InfluxdbService{Url: cf.Influxdb.Url}

	sp := &Container{
		Config:   cf,
		Redis:    redisService,
		Mongo:    mongo,
		Influxdb: influx,
	}

	return sp
}

func NewContainer(configPath string) *Container {
	cf := config.MustRead(configPath)
	return New(cf)
}
