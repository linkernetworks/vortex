package prometheus

import (
	"github.com/linkernetworks/logger"
	"github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
)

type PrometheusConfig struct {
	Url string `json:"url"`
}

type Service struct {
	Url string
	API prometheus.API
}

func New(url string) *Service {

	conf := api.Config{
		Address:      url,
		RoundTripper: api.DefaultRoundTripper,
	}

	client, err := api.NewClient(conf)
	if err != nil {
		logger.Fatalf("error while creating api.NewClient %s", err)
	}

	newAPI := prometheus.NewAPI(client)

	return &Service{
		Url: url,
		API: newAPI,
	}
}
