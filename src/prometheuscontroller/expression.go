package prometheuscontroller

import (
	"strings"

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

func GetNode(sp *serviceprovider.Container, id string) (entity.NodeMetrics, error) {
	node := entity.NodeMetrics{}
	node.Detail.Labels = map[string]string{}
	node.NICs = map[string]entity.NICMetrics{}

	// kube_node_info, kube_node_created, node_network_interface, kube_node_labels, kube_node_status_capacity, kube_node_status_allocatable
	results, err := query(sp, `{__name__=~"kube_node_info|kube_node_created|node_network_interface|kube_node_labels|kube_node_status_capacity|kube_node_status_allocatable",node=~"`+id+`"}`)
	if err != nil {
		return node, err
	}

	for _, result := range results {
		switch result.Metric["__name__"] {

		case "kube_node_info":
			node.Detail.Hostname = id
			node.Detail.KernelVersion = string(result.Metric["kernel_version"])
			node.Detail.KubeproxyVersion = string(result.Metric["kubeproxy_version"])
			node.Detail.OS = string(result.Metric["os_image"])
			node.Detail.KubernetesVersion = string(result.Metric["kubelet_version"])

		case "kube_node_created":
			node.Detail.CreatedAt = int(result.Value)

		case "kube_node_labels":
			for key, value := range result.Metric {
				if strings.HasPrefix(string(key), "label_") {
					node.Detail.Labels[strings.TrimPrefix(string(key), "label_")] = string(value)
				}
			}

		case "kube_node_status_allocatable":
			switch result.Metric["resource"] {
			case "cpu":
				node.Resource.AllocatableCPU = float32(result.Value)
			case "memory":
				node.Resource.AllocatableMemory = float32(result.Value)
			case "pods":
				node.Resource.AllocatablePods = float32(result.Value)
			case "ephemeral_storage":
				node.Resource.AllocatableEphemeralStorage = float32(result.Value)
			}

		case "kube_node_status_capacity":
			switch result.Metric["resource"] {
			case "cpu":
				node.Resource.CapacityCPU = float32(result.Value)
			case "memory":
				node.Resource.CapacityMemory = float32(result.Value)
			case "pods":
				node.Resource.CapacityPods = float32(result.Value)
			case "ephemeral_storage":
				node.Resource.CapacityEphemeralStorage = float32(result.Value)
			}
		}
	}

	// kube_node_status_condition
	results, err = query(sp, `{__name__=~"kube_node_status_condition",node=~"`+id+`",status="true"}==1`)
	if err != nil {
		return node, err
	}

	node.Detail.Status = string(results[0].Metric["condition"])

	// kube_pod_container_resource_limits, kube_pod_container_resource_requests
	results, err = query(sp, `sum by(__name__, resource) ({__name__=~"kube_pod_container_resource_limits|kube_pod_container_resource_requests",node=~"`+id+`"})`)
	if err != nil {
		return node, err
	}

	for _, result := range results {
		switch result.Metric["__name__"] {
		case "kube_pod_container_resource_requests":
			switch result.Metric["resource"] {
			case "cpu":
				node.Resource.CPURequests = float32(result.Value)
			case "memory":
				node.Resource.MemoryRequests = float32(result.Value)
			}
		case "kube_pod_container_resource_limits":
			switch result.Metric["resource"] {
			case "cpu":
				node.Resource.CPULimits = float32(result.Value)
			case "memory":
				node.Resource.MemoryLimits = float32(result.Value)
			}
		}
	}

	// node_network_interface, node_network_receive_bytes_total, node_network_transmit_bytes_total, node_network_receive_packets_total, node_network_transmit_packets_total
	results, err = query(sp, `{__name__=~"node_network_interface|node_network_receive_bytes_total|node_network_transmit_bytes_total|node_network_receive_packets_total|node_network_transmit_packets_total",node=~"`+id+`"}`)
	if err != nil {
		return node, err
	}

	for _, result := range results {
		switch result.Metric["__name__"] {

		case "node_network_interface":
			nic := entity.NICMetrics{}
			nic.Default = string(result.Metric["default"])
			nic.Type = string(result.Metric["type"])
			nic.IP = string(result.Metric["ip_address"])
			nic.NICNetworkTraffic = entity.NICNetworkTrafficMetrics{}
			node.NICs[string(result.Metric["device"])] = nic

		case "node_network_receive_bytes_total":
			nic := node.NICs[string(result.Metric["device"])]
			nic.NICNetworkTraffic.ReceiveBytesTotal = int(result.Value)
			node.NICs[string(result.Metric["device"])] = nic

		case "node_network_transmit_bytes_total":
			nic := node.NICs[string(result.Metric["device"])]
			nic.NICNetworkTraffic.TransmitBytesTotal = int(result.Value)
			node.NICs[string(result.Metric["device"])] = nic

		case "node_network_receive_packets_total":
			nic := node.NICs[string(result.Metric["device"])]
			nic.NICNetworkTraffic.ReceivePacketsTotal = int(result.Value)
			node.NICs[string(result.Metric["device"])] = nic

		case "node_network_transmit_packets_total":
			nic := node.NICs[string(result.Metric["device"])]
			nic.NICNetworkTraffic.TransmitPacketsTotal = int(result.Value)
			node.NICs[string(result.Metric["device"])] = nic
		}
	}

	return node, nil
}
