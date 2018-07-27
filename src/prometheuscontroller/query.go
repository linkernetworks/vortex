package prometheuscontroller

import (
	"fmt"
	"strings"
	"time"

	"github.com/linkernetworks/vortex/src/serviceprovider"
	pv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"golang.org/x/net/context"
)

// Expression is the structure for expression
type Expression struct {
	Metrics     []string          `json:"metrics"`
	QueryLabels map[string]string `json:"queryLabels"`
	SumBy       []string          `json:"sumBy"`
	Value       *int              `json:"value"`
	Time        *string           `json:"time"`
}

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

func basicExpr(expression Expression) (string, error) {
	var str string

	// append the metrics
	var metrics string
	for _, metric := range expression.Metrics {
		metrics = fmt.Sprintf("%s%s|", metrics, metric)
	}
	str = fmt.Sprintf(`__name__=~"%s"`, strings.TrimSuffix(metrics, "|"))

	// add the query labels
	var labels string
	for key, value := range expression.QueryLabels {
		labels = fmt.Sprintf(`,%s=~"%s"`, key, value)
	}
	str = fmt.Sprintf("{%s%s}", str, labels)

	// use sum by if need it
	if expression.SumBy != nil {
		var sumby string
		for _, sumLabel := range expression.SumBy {
			sumby = fmt.Sprintf("%s%s,", sumby, sumLabel)
		}
		str = fmt.Sprintf("sum by(%s)(%s)", strings.TrimSuffix(sumby, ","), str)
	}

	// durations over the last [] time
	if expression.Time != nil {
		str = fmt.Sprintf("%s[%v]", str, *expression.Time)
	}

	// the result should equal to expression.Value
	if expression.Value != nil {
		str = fmt.Sprintf("%s==%v", str, *expression.Value)
	}

	return str, nil

	// results, err := query(sp, str)
	// if err != nil {
	// 	return nil, fmt.Errorf("%v, can not query the expression: %s", err, str)
	// }
	// return results, nil
}
