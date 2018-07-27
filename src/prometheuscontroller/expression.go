package prometheuscontroller

import (
	"fmt"
	"strings"
	"time"
)

// Expression is the structure for expression
type Expression struct {
	Metrics     []string          `json:"metrics"`
	QueryLabels map[string]string `json:"queryLabels"`
	SumByLabels []string          `json:"sumByLabels"`
}

// Create basic expression with metrics
func basicExpr(metrics []string) (string, error) {
	if metrics == nil {
		return "", fmt.Errorf("no metrics be assigned")
	}

	var tmp string
	for _, metric := range metrics {
		tmp = fmt.Sprintf("%s%s|", tmp, metric)
	}
	expr := fmt.Sprintf(`__name__=~"%s"`, strings.TrimSuffix(tmp, "|"))

	return expr, nil
}

// Append the query labels for expression
func queryExpr(expr string, queryLabels map[string]string) string {
	if queryLabels == nil {
		return expr
	}

	var tmp string
	for key, value := range queryLabels {
		tmp = fmt.Sprintf(`,%s=~"%s"`, key, value)
	}
	expr = fmt.Sprintf("{%s%s}", expr, tmp)

	return expr
}

// Append the sum by syntax with labels
func sumByExpr(expr string, sumByLabels []string) string {
	if sumByLabels == nil {
		return expr
	}

	var tmp string
	for _, sumLabel := range sumByLabels {
		tmp = fmt.Sprintf("%s%s,", tmp, sumLabel)
	}
	expr = fmt.Sprintf("sum by(%s)(%s)", strings.TrimSuffix(tmp, ","), expr)

	return expr
}

// Append the sum syntax
func sumExpr(expr string) string {
	expr = fmt.Sprintf("sum(%s)", expr)

	return expr
}

// Append a duration for expression
func durationExpr(expr string, duration time.Duration) string {
	expr = fmt.Sprintf("%s[%v]", expr, duration)

	return expr
}

// Append the rate syntax
func rateExpr(expr string) string {
	expr = fmt.Sprintf("rate(%s)", expr)

	return expr
}

// Assign a value for expression
func equalExpr(expr string, value float64) string {
	expr = fmt.Sprintf("%s==%v", expr, value)

	return expr
}
