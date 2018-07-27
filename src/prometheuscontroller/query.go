package prometheuscontroller

import (
	"fmt"
	"time"

	"github.com/linkernetworks/vortex/src/serviceprovider"
	pv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"golang.org/x/net/context"
)

func query(sp *serviceprovider.Container, expression string) (model.Vector, error) {
	api := sp.Prometheus.API

	testTime := time.Now()
	result, err := api.Query(context.Background(), expression, testTime)

	// https://github.com/prometheus/client_golang/blob/d6a9817c4afc94d51115e4a30d449056a3fbf547/api/prometheus/v1/api.go#L316
	// this api always return the err no matter what
	// so we should use result==nil to determine whether it is a true error
	if result == nil {
		return nil, err
	}

	if result.Type() == model.ValVector {
		return result.(model.Vector), nil
	}
	return nil, fmt.Errorf("the type of the return result can not be identify")
}

func queryRange(sp *serviceprovider.Container, expression string) (model.Matrix, error) {
	api := sp.Prometheus.API

	rangeSet := pv1.Range{Start: time.Now().Add(-time.Minute * 2), End: time.Now(), Step: time.Second * 10}
	result, err := api.QueryRange(context.Background(), expression, rangeSet)

	// https://github.com/prometheus/client_golang/blob/d6a9817c4afc94d51115e4a30d449056a3fbf547/api/prometheus/v1/api.go#L316
	// this api always return the err no matter what
	// so we should use result==nil to determine whether it is a true error
	if result == nil {
		return nil, err
	}

	if result.Type() == model.ValMatrix {
		return result.(model.Matrix), nil
	}
	return nil, fmt.Errorf("the type of the return result can not be identify")
}
