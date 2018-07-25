package prometheuscontroller

import (
	"strconv"
	"strings"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/prometheus/common/model"
)

// ListResource will list resource
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

// ListNodeNICs will list node's NICs
func ListNodeNICs(sp *serviceprovider.Container, id string) (entity.NodeNICsMetrics, error) {
	nicList := entity.NodeNICsMetrics{}
	expression := Expression{}
	expression.Metrics = []string{"node_network_interface"}
	expression.QueryLabels = map[string]string{"node": id}

	results, err := getElements(sp, expression)
	if err != nil {
		return nicList, err
	}

	for _, result := range results {
		nic := entity.NICOverviewMetrics{}
		nic.Name = string(result.Metric["device"])
		nic.Type = string(result.Metric["type"])
		nic.PCIID = string(result.Metric["pci_id"])
		defaultValue, err := strconv.ParseBool(string(result.Metric["default"]))
		if err != nil {
			return nicList, err
		}
		nic.Default = defaultValue

		nicList.NICs = append(nicList.NICs, nic)
	}

	return nicList, nil
}

// GetPod will get pod
func GetPod(sp *serviceprovider.Container, id string) (entity.PodMetrics, error) {
	pod := entity.PodMetrics{}
	pod.Labels = map[string]string{}
	pod.NICs = map[string]entity.NICShortMetrics{}

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

	// network interface
	expression = Expression{}
	expression.Metrics = []string{"container_network_receive_bytes_total"}
	expression.QueryLabels = map[string]string{"container_label_io_kubernetes_pod_name": id}

	results, err = getElements(sp, expression)
	if err != nil {
		return pod, err
	}

	for _, result := range results {
		nic := entity.NICShortMetrics{}
		nic.NICNetworkTraffic = entity.NICNetworkTrafficMetrics{}
		pod.NICs[string(result.Metric["interface"])] = nic
	}

	// network traffic
	expression = Expression{}
	expression.Metrics = []string{"container_network_receive_bytes_total", "container_network_transmit_bytes_total", "container_network_receive_packets_total", "container_network_transmit_packets_total"}
	expression.QueryLabels = map[string]string{"container_label_io_kubernetes_pod_name": id}

	results, err = getElements(sp, expression)
	if err != nil {
		return pod, err
	}

	for _, result := range results {
		switch result.Metric["__name__"] {

		case "container_network_receive_bytes_total":
			nic := pod.NICs[string(result.Metric["interface"])]
			nic.NICNetworkTraffic.ReceiveBytesTotal = int(result.Value)
			pod.NICs[string(result.Metric["interface"])] = nic

		case "container_network_transmit_bytes_total":
			nic := pod.NICs[string(result.Metric["interface"])]
			nic.NICNetworkTraffic.TransmitBytesTotal = int(result.Value)
			pod.NICs[string(result.Metric["interface"])] = nic

		case "container_network_receive_packets_total":
			nic := pod.NICs[string(result.Metric["interface"])]
			nic.NICNetworkTraffic.ReceivePacketsTotal = int(result.Value)
			pod.NICs[string(result.Metric["interface"])] = nic

		case "container_network_transmit_packets_total":
			nic := pod.NICs[string(result.Metric["interface"])]
			nic.NICNetworkTraffic.TransmitPacketsTotal = int(result.Value)
			pod.NICs[string(result.Metric["interface"])] = nic

		}
	}

	return pod, nil
}

// GetContainer will get container
func GetContainer(sp *serviceprovider.Container, id string) (entity.ContainerMetrics, error) {
	container := entity.ContainerMetrics{}

	// basic info
	expression := Expression{}
	expression.Metrics = []string{"kube_pod_container_info", "kube_pod_container_status_restarts_total"}
	expression.QueryLabels = map[string]string{"container": id}

	results, err := getElements(sp, expression)
	if err != nil {
		return container, err
	}

	for _, result := range results {
		switch result.Metric["__name__"] {

		case "kube_pod_container_info":
			container.Detail.ContainerName = id
			container.Detail.Pod = string(result.Metric["pod"])
			container.Detail.Node = string(result.Metric["node"])
			container.Detail.Image = string(result.Metric["image"])
			container.Detail.Namespace = string(result.Metric["namespace"])

		case "kube_pod_container_status_restarts_total":
			container.Status.RestartTime = int(result.Value)
		}
	}

	// status
	expression = Expression{}
	expression.Metrics = []string{"kube_pod_container_status.*"}
	expression.QueryLabels = map[string]string{"container": id}
	ValueInt := 1
	expression.Value = &ValueInt

	results, err = getElements(sp, expression)
	if err != nil {
		return container, err
	}

	for _, result := range results {
		switch result.Metric["__name__"] {

		case "kube_pod_container_status_ready":
			if container.Status.Status == "" {
				container.Status.Status = "ready"
			}

		case "kube_pod_container_status_running":
			container.Status.Status = "running"

		case "kube_pod_container_status_waiting":
			container.Status.Status = "waiting"

		case "kube_pod_container_status_terminated":
			container.Status.Status = "terminated"

		case "kube_pod_container_status_terminated_reason":
			container.Status.TerminatedReason = string(result.Metric["reason"])

		case "kube_pod_container_status_waiting_reason":
			container.Status.WaitingReason = string(result.Metric["reason"])
		}
	}

	expression = Expression{}
	expression.Metrics = []string{"container_last_seen"}
	expression.QueryLabels = map[string]string{"container_label_io_kubernetes_container_name": id}
	results, err = getElements(sp, expression)
	if err != nil {
		return container, err
	}
	if len(results) == 0 {
		return container, err
	}

	// resource
	results, err = query(sp, `sum(rate(container_cpu_usage_seconds_total{container_label_io_kubernetes_container_name=~"`+id+`"}[1m])) * 100`)
	if err != nil {
		return container, err
	}
	container.Resource.CPUUsagePercentage = float32(results[0].Value)

	expression = Expression{}
	expression.Metrics = []string{"container_memory_usage_bytes"}
	expression.QueryLabels = map[string]string{"container_label_io_kubernetes_container_name": id}

	results, err = getElements(sp, expression)
	if err != nil {
		return container, err
	}

	container.Resource.MemoryUsageBytes = float32(results[0].Value)

	// command
	kc := sp.KubeCtl
	pod, err := kc.GetPod(container.Detail.Pod, container.Detail.Namespace)
	if err != nil {
		return entity.ContainerMetrics{}, err
	}

	for _, obj := range pod.Spec.Containers {
		if obj.Name == id {
			container.Detail.Command = obj.Command
			break
		}
	}

	return container, nil
}

// GetService will get service by serviceprovider.Container
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

	// get service port config
	kc := sp.KubeCtl
	object, err := kc.GetService(service.ServiceName, service.Namespace)
	if err != nil {
		return entity.ServiceMetrics{}, err
	}
	service.Ports = object.Spec.Ports

	return service, nil
}

// GetController willl get container
func GetController(sp *serviceprovider.Container, id string) (entity.ControllerMetrics, error) {
	controller := entity.ControllerMetrics{}
	controller.Labels = map[string]string{}
	controller.Type = "deployment"

	expression := Expression{}
	expression.Metrics = []string{"kube_deployment_metadata_generation", "kube_deployment_created", "kube_deployment_labels", "kube_deployment_spec_replicas", "kube_deployment_status_replicas", "kube_deployment_status_replicas_available"}
	expression.QueryLabels = map[string]string{"deployment": id}

	results, err := getElements(sp, expression)
	if err != nil {
		return controller, err
	}

	for _, result := range results {
		switch result.Metric["__name__"] {

		case "kube_deployment_metadata_generation":
			controller.ControllerName = id
			controller.Namespace = string(result.Metric["namespace"])

		case "kube_deployment_spec_replicas":
			controller.DesiredPod = int(result.Value)

		case "kube_deployment_status_replicas":
			controller.CurrentPod = int(result.Value)

		case "kube_deployment_status_replicas_available":
			controller.AvailablePod = int(result.Value)

		case "kube_deployment_created":
			controller.CreateAt = int(result.Value)

		case "kube_deployment_labels":
			for key, value := range result.Metric {
				if strings.HasPrefix(string(key), "label_") {
					controller.Labels[strings.TrimPrefix(string(key), "label_")] = string(value)
				}
			}
		}

	}

	return controller, nil
}

// GetNode will get node metrics
func GetNode(sp *serviceprovider.Container, id string) (entity.NodeMetrics, error) {
	node := entity.NodeMetrics{}
	node.Detail.Labels = map[string]string{}
	node.NICs = map[string]entity.NICMetrics{}

	// basic info
	expression := Expression{}
	expression.Metrics = []string{"kube_node_info", "kube_node_created", "node_network_interface", "kube_node_labels", "kube_node_status_capacity", "kube_node_status_allocatable"}
	expression.QueryLabels = map[string]string{"node": id}

	results, err := getElements(sp, expression)
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

		case "node_network_interface":
			nic := entity.NICMetrics{}
			defaultValue, err := strconv.ParseBool(string(result.Metric["default"]))
			if err != nil {
				return node, err
			}
			nic.Default = defaultValue
			nic.Type = string(result.Metric["type"])
			nic.IP = string(result.Metric["ip_address"])
			nic.PCIID = string(result.Metric["pci_id"])
			nic.NICNetworkTraffic = entity.NICNetworkTrafficMetrics{}
			node.NICs[string(result.Metric["device"])] = nic

		case "kube_node_status_allocatable":
			switch result.Metric["resource"] {
			case "cpu":
				node.Resource.AllocatableCPU = float32(result.Value)
			case "memory":
				node.Resource.AllocatableMemory = float32(result.Value)
			case "pods":
				node.Resource.AllocatablePods = float32(result.Value)
			}

		case "kube_node_status_capacity":
			switch result.Metric["resource"] {
			case "cpu":
				node.Resource.CapacityCPU = float32(result.Value)
			case "memory":
				node.Resource.CapacityMemory = float32(result.Value)
			case "pods":
				node.Resource.CapacityPods = float32(result.Value)
			}
		}
	}

	// status
	expression = Expression{}
	expression.Metrics = []string{"kube_node_status_condition"}
	expression.QueryLabels = map[string]string{"node": id, "status": "true"}
	ValueInt := 1
	expression.Value = &ValueInt

	results, err = getElements(sp, expression)
	if err != nil {
		return node, err
	}

	node.Detail.Status = string(results[0].Metric["condition"])

	// resource
	expression = Expression{}
	expression.Metrics = []string{"kube_pod_container_resource_limits", "kube_pod_container_resource_requests"}
	expression.QueryLabels = map[string]string{"node": id}
	expression.SumBy = []string{"__name__", "resource"}

	results, err = getElements(sp, expression)
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

	// network traffic
	expression = Expression{}
	expression.Metrics = []string{"node_network_receive_bytes_total", "node_network_transmit_bytes_total", "node_network_receive_packets_total", "node_network_transmit_packets_total"}
	expression.QueryLabels = map[string]string{"node": id}

	results, err = getElements(sp, expression)
	if err != nil {
		return node, err
	}

	for _, result := range results {
		switch result.Metric["__name__"] {

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
