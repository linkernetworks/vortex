package prometheuscontroller

import (
	"fmt"
	"strings"
	"time"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/vortex/src/serviceprovider"
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

	switch {
	case result.Type() == model.ValVector:
		return result.(model.Vector), nil
	default:
		return nil, fmt.Errorf("the type of the return result can not be identify")
	}
}

func getElements(sp *serviceprovider.Container, expression Expression) (model.Vector, error) {
	// append the metrics
	str := `__name__=~"{{metrics}}"`
	metrics := ""
	for _, metric := range expression.Metrics {
		metrics = metrics + metric + "|"
	}
	rule := strings.NewReplacer("{{metrics}}", strings.TrimSuffix(metrics, "|"))
	str = rule.Replace(str)

	// add the query labels
	labels := expression.QueryLabels
	for key, value := range labels {
		str = fmt.Sprintf(`%s,%s=~"%s"`, str, key, value)
	}
	str = `{` + str + `}`

	// use sum by if need it
	if expression.SumBy != nil {
		str = fmt.Sprintf(`sum by({{sumby}})(%s)`, str)
		sumby := ""
		for _, sumLabel := range expression.SumBy {
			sumby = sumby + sumLabel + ","
		}
		rule = strings.NewReplacer("{{sumby}}", strings.TrimSuffix(sumby, ","))
		str = rule.Replace(str)
	}

	if expression.Value != nil {
		str = fmt.Sprintf(`%s==%v`, str, *expression.Value)
	}

	logger.Infof(str)

	results, err := query(sp, str)
	if err != nil {
		return nil, fmt.Errorf("%v, can not query the expression: %s", err, str)
	}
	return results, nil
}
