package prometheuscontroller

import (
	"fmt"
	"strings"
	"time"

	"github.com/linkernetworks/vortex/src/serviceprovider"
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

func query(sp *serviceprovider.Container, expression string) (interface{}, error) {

	api := sp.Prometheus.API

	testTime := time.Now()
	result, err := api.Query(context.Background(), expression, testTime)

	// https://github.com/prometheus/client_golang/blob/d6a9817c4afc94d51115e4a30d449056a3fbf547/api/prometheus/v1/api.go#L316
	// this api always return the err no matter what
	// so we should use result==nil to determine whether it is a true error
	if result == nil {
		return nil, err
	}

	if result.Type() == model.ValVector || result.Type() == model.ValMatrix {
		return result, nil
	} else {
		return nil, fmt.Errorf("the type of the return result can not be identify")
	}

}

func getElements(sp *serviceprovider.Container, expression Expression) (interface{}, error) {
	// append the metrics
	var metrics string
	str := `__name__=~"{{metrics}}"`
	for _, metric := range expression.Metrics {
		metrics = fmt.Sprintf("%s%s|", metrics, metric)
	}
	rule := strings.NewReplacer("{{metrics}}", strings.TrimSuffix(metrics, "|"))
	str = rule.Replace(str)

	// add the query labels
	labels := expression.QueryLabels
	for key, value := range labels {
		str = fmt.Sprintf(`%s,%s=~"%s"`, str, key, value)
	}
	str = fmt.Sprintf("{%s}", str)

	// use sum by if need it
	var sumby string
	if expression.SumBy != nil {
		str = fmt.Sprintf("sum by({{sumby}})(%s)", str)
		for _, sumLabel := range expression.SumBy {
			sumby = fmt.Sprintf("%s%s,", sumby, sumLabel)
		}
		rule = strings.NewReplacer("{{sumby}}", strings.TrimSuffix(sumby, ","))
		str = rule.Replace(str)
	}

	// the result should equal to expression.Value
	if expression.Value != nil {
		str = fmt.Sprintf("%s==%v", str, *expression.Value)
	}

	// return the historical result
	if expression.Time != nil {
		str = fmt.Sprintf("%s[%v]", str, *expression.Time)
	}

	results, err := query(sp, str)
	if err != nil {
		return nil, fmt.Errorf("%v, can not query the expression: %s", err, str)
	}
	return results, nil
}
