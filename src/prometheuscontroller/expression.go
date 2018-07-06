package prometheuscontroller

import (
	"strings"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/prometheus/common/model"
)

func ListResource(sp *serviceprovider.Container, resource model.LabelName, mapList map[string]string) ([]string, error) {
	resourceList := []string{}
	expression := `({__name__=~"{{metric}}",node=~"{{node}}", namespace=~"{{namespace}}"})`
	expression = replacer(expression, mapList)

	logger.Infof(expression)

	results, err := query(sp, expression)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		resourceList = append(resourceList, string(result.Metric[resource]))
	}

	return resourceList, nil
}

func replacer(str string, mapList map[string]string) string {
	for key, value := range mapList {
		key = "{{" + key + "}}"
		str = strings.Replace(str, key, value, -1)
	}
	return str
}
