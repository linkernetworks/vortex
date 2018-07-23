package prometheusprovider

import (
	"github.com/linkernetworks/logger"
	"github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
)

// PrometheusConfig is the structure for Prometheus Config
type PrometheusConfig struct {
	Url string `json:"url"`
}

// Service is the structure for Service
type Service struct {
	Url string
	API prometheus.API
}

// New will reture a new service
func New(url string) *Service {
	conf := api.Config{
		Address:      url,
		RoundTripper: api.DefaultRoundTripper,
	}

	client, err := api.NewClient(conf)
	if err != nil {
		// TODO should return error to server
		logger.Warnf("error while creating api.NewClient %s", err)
	}

	newAPI := prometheus.NewAPI(client)

	return &Service{
		Url: url,
		API: newAPI,
	}
}
