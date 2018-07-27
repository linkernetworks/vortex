package prometheuscontroller

import (
	"fmt"
	"strings"
)

// Expression is the structure for expression
type Expression struct {
	Metrics     []string          `json:"metrics"`
	QueryLabels map[string]string `json:"queryLabels"`
	SumBy       []string          `json:"sumBy"`
	Value       *int              `json:"value"`
	Time        *string           `json:"time"`
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
