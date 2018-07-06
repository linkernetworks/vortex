package prometheuscontroller

import (
	"fmt"
	"strings"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/prometheus/common/model"
)

type Expression struct {
	Metrics     []string          `json:"metrics"`
	QueryLabels map[string]string `json:"queryLabels"`
	SumBy       []string          `json:"sumBy"`
	Value       *int              `json:"value"`
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

func ListResource(sp *serviceprovider.Container, resource model.LabelName, expression Expression) ([]string, error) {
	results, err := getElements(sp, expression)
	if err != nil {
		return nil, err
	}

	resourceList := []string{}
	for _, result := range results {
		resourceList = append(resourceList, string(result.Metric[resource]))
	}

	return resourceList, nil
}

func GetPod(sp *serviceprovider.Container, id string) (entity.PodMetrics, error) {
	pod := entity.PodMetrics{}
	pod.Labels = map[string]string{}

	expression := Expression{}
	expression.Metrics = []string{"kube_pod_info", "kube_pod_created", "kube_pod_labels", "kube_pod_owner", "kube_pod_status_phase", "kube_pod_container_info", "kube_pod_container_status_restarts_total"}
	expression.QueryLabels = map[string]string{"pod": id}

	results, err := getElements(sp, expression)
	if err != nil {
		return pod, err
	}

	for _, result := range results {
		switch result.Metric["__name__"] {

		case "kube_pod_info":
			pod.PodName = id
			pod.IP = string(result.Metric["pod_ip"])
			pod.Node = string(result.Metric["node"])
			pod.Namespace = string(result.Metric["namespace"])
			pod.CreateByKind = string(result.Metric["created_by_kind"])
			pod.CreateByName = string(result.Metric["created_by_name"])

		case "kube_pod_created":
			pod.CreateAt = int(result.Value)

		case "kube_pod_labels":
			for key, value := range result.Metric {
				if strings.HasPrefix(string(key), "label_") {
					pod.Labels[strings.TrimPrefix(string(key), "label_")] = string(value)
				}
			}

		case "kube_pod_status_phase":
			if int(result.Value) == 1 {
				pod.Status = string(result.Metric["phase"])
			}

		case "kube_pod_container_info":
			pod.Containers = append(pod.Containers, string(result.Metric["container"]))

		case "kube_pod_container_status_restarts_total":
			pod.RestartCount = pod.RestartCount + int(result.Value)
		}
	}

	return pod, nil
}

func GetService(sp *serviceprovider.Container, id string) (entity.ServiceMetrics, error) {
	service := entity.ServiceMetrics{}
	service.Labels = map[string]string{}

	expression := Expression{}
	expression.Metrics = []string{"kube_service_info", "kube_service_created", "kube_service_labels", "kube_service_spec_type"}
	expression.QueryLabels = map[string]string{"service": id}

	results, err := getElements(sp, expression)
	if err != nil {
		return service, err
	}

	for _, result := range results {
		switch result.Metric["__name__"] {

		case "kube_service_info":
			service.ServiceName = id
			service.Namespace = string(result.Metric["namespace"])
			service.ClusterIP = string(result.Metric["cluster_ip"])

		case "kube_service_spec_type":
			service.Type = string(result.Metric["type"])

		case "kube_service_created":
			service.CreateAt = int(result.Value)

		case "kube_service_labels":
			for key, value := range result.Metric {
				if strings.HasPrefix(string(key), "label_") {
					service.Labels[strings.TrimPrefix(string(key), "label_")] = string(value)
				}
			}
		}

	}

	return service, nil
}
