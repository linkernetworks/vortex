package prometheuscontroller

import (
	"fmt"
	"strings"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/prometheus/common/model"
)

type Expression struct {
	Metrics      []string          `json:"metrics"`
	QueryLabels  map[string]string `json:"queryLabels"`
	TargetLabels []model.LabelName `json:"targetLabels"`
	SumBy        []string          `json:"sumBy"`
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

	logger.Infof(str)

	results, err := query(sp, str)
	if err != nil {
		return nil, fmt.Errorf("%v, can not query the expression: %s", err, str)
	}
	return results, nil
}

func ListResource(sp *serviceprovider.Container, expression Expression) ([]string, error) {
	resourceList := []string{}
	results, err := getElements(sp, expression)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		resourceList = append(resourceList, string(result.Metric[expression.TargetLabels[0]]))
	}

	return resourceList, nil
}
