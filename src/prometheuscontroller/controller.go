package prometheuscontroller

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

// ListContainerName will list container name
func ListContainerName(sp *serviceprovider.Container, queryLabels map[string]string) ([]string, error) {
	expression := Expression{}
	expression.Metrics = []string{"kube_pod_container_info"}
	expression.QueryLabels = queryLabels

	str := basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	results, err := query(sp, str)
	if err != nil {
		return nil, err
	}

	containerList := []string{}
	for _, result := range results {
		containerList = append(containerList, string(result.Metric["container"]))
	}

	return containerList, nil
}

// ListPodName will list pod name
func ListPodName(sp *serviceprovider.Container, queryLabels map[string]string) ([]string, error) {
	expression := Expression{}
	expression.Metrics = []string{"kube_pod_info"}
	expression.QueryLabels = queryLabels

	str := basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	results, err := query(sp, str)
	if err != nil {
		return nil, err
	}

	podList := []string{}
	for _, result := range results {
		podList = append(podList, string(result.Metric["pod"]))
	}

	return podList, nil
}

// ListServiceName will list service name
func ListServiceName(sp *serviceprovider.Container, queryLabels map[string]string) ([]string, error) {
	expression := Expression{}
	expression.Metrics = []string{"kube_service_info"}
	expression.QueryLabels = queryLabels

	str := basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	results, err := query(sp, str)
	if err != nil {
		return nil, err
	}

	serviceList := []string{}
	for _, result := range results {
		serviceList = append(serviceList, string(result.Metric["service"]))
	}

	return serviceList, nil
}

// ListControllerName will list controller name
func ListControllerName(sp *serviceprovider.Container, queryLabels map[string]string) ([]string, error) {
	expression := Expression{}
	expression.Metrics = []string{"kube_deployment_metadata_generation"}
	expression.QueryLabels = queryLabels

	str := basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	results, err := query(sp, str)
	if err != nil {
		return nil, err
	}

	controllerList := []string{}
	for _, result := range results {
		controllerList = append(controllerList, string(result.Metric["deployment"]))
	}

	return controllerList, nil
}

// ListNodeName will list node name
func ListNodeName(sp *serviceprovider.Container, queryLabels map[string]string) ([]string, error) {
	expression := Expression{}
	expression.Metrics = []string{"kube_node_info"}
	expression.QueryLabels = queryLabels

	str := basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	results, err := query(sp, str)
	if err != nil {
		return nil, err
	}

	nodeList := []string{}
	for _, result := range results {
		nodeList = append(nodeList, string(result.Metric["exported_node"]))
	}

	return nodeList, nil
}

// ListNodeNICs will list node's NICs
func ListNodeNICs(sp *serviceprovider.Container, id string) (entity.NodeNICsMetrics, error) {
	nicList := entity.NodeNICsMetrics{}
	expression := Expression{}
	expression.Metrics = []string{"node_network_interface"}
	expression.QueryLabels = map[string]string{"node": id}

	str := basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	results, err := query(sp, str)
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
		dpdkValue, err := strconv.ParseBool(string(result.Metric["dpdk"]))
		if err != nil {
			return nicList, err
		}
		nic.DPDK = dpdkValue

		nicList.NICs = append(nicList.NICs, nic)
	}

	return nicList, nil
}

// GetPod will get pod
func GetPod(sp *serviceprovider.Container, id string, rs RangeSetting) (entity.PodMetrics, error) {
	pod := entity.PodMetrics{}
	pod.Labels = map[string]string{}
	pod.NICs = map[string]entity.NICShortMetrics{}

	expression := Expression{}
	expression.Metrics = []string{
		"kube_pod_info",
		"kube_pod_created",
		"kube_pod_labels",
		"kube_pod_owner",
		"kube_pod_status_phase",
		"kube_pod_container_info",
		"kube_pod_container_status_restarts_total"}
	expression.QueryLabels = map[string]string{"pod": id}

	str := basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	results, err := query(sp, str)
	if err != nil {
		return pod, err
	}

	for _, result := range results {
		switch result.Metric["__name__"] {

		case "kube_pod_info":
			pod.PodName = id
			pod.IP = string(result.Metric["pod_ip"])
			pod.Node = string(result.Metric["exported_node"])
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
				pod.Phase = string(result.Metric["phase"])
			}

		case "kube_pod_container_info":
			pod.Containers = append(pod.Containers, string(result.Metric["container"]))

		case "kube_pod_container_status_restarts_total":
			pod.RestartCount = pod.RestartCount + int(result.Value)
		}
	}

	// status
	pod.Status = ""
podStatusCheckingLoop:
	for _, container := range pod.Containers {
		expression.Metrics = []string{"kube_pod_container_status.*"}
		expression.QueryLabels = map[string]string{"container": container}

		str = basicExpr(expression.Metrics)
		str = queryExpr(str, expression.QueryLabels)
		str = equalExpr(str, 1)

		results, err = query(sp, str)
		if err != nil {
			return pod, err
		}

		for _, result := range results {
			switch result.Metric["__name__"] {

			case "kube_pod_container_status_ready":
				pod.Status = "ready"

			case "kube_pod_container_status_running":
				pod.Status = "running"

			case "kube_pod_container_status_terminated_reason":
				pod.Status = string(result.Metric["reason"])
				break podStatusCheckingLoop

			case "kube_pod_container_status_waiting_reason":
				pod.Status = string(result.Metric["reason"])
				break podStatusCheckingLoop
			}
		}
	}

	// network interface
	expression.Metrics = []string{"container_network_receive_bytes_total"}
	expression.QueryLabels = map[string]string{"container_label_io_kubernetes_pod_name": id}

	str = basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	results, err = query(sp, str)
	if err != nil {
		return pod, err
	}

	for _, result := range results {
		nic := entity.NICShortMetrics{}
		nic.NICNetworkTraffic = entity.NICNetworkTrafficMetrics{}
		pod.NICs[string(result.Metric["interface"])] = nic
	}

	// network traffic receive bytes
	expression.Metrics = []string{"container_network_receive_bytes_total"}
	expression.QueryLabels = map[string]string{"container_label_io_kubernetes_pod_name": id}

	str = basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	str = rateExpr(durationExpr(str, int(rs.Rate)))
	resultMatrix, err := queryRange(sp, str, rs)
	if err != nil {
		return pod, err
	}

	for _, result := range resultMatrix {
		nic := pod.NICs[string(result.Metric["interface"])]
		for _, pair := range result.Values {
			nic.NICNetworkTraffic.ReceiveBytesTotal = append(nic.NICNetworkTraffic.ReceiveBytesTotal, entity.SamplePair{Timestamp: pair.Timestamp, Value: pair.Value})
		}
		pod.NICs[string(result.Metric["interface"])] = nic
	}

	// network traffic transmit bytes
	expression.Metrics = []string{"container_network_transmit_bytes_total"}
	expression.QueryLabels = map[string]string{"container_label_io_kubernetes_pod_name": id}

	str = basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	str = rateExpr(durationExpr(str, int(rs.Rate)))
	resultMatrix, err = queryRange(sp, str, rs)
	if err != nil {
		return pod, err
	}

	for _, result := range resultMatrix {
		nic := pod.NICs[string(result.Metric["interface"])]
		for _, pair := range result.Values {
			nic.NICNetworkTraffic.TransmitBytesTotal = append(nic.NICNetworkTraffic.TransmitBytesTotal, entity.SamplePair{Timestamp: pair.Timestamp, Value: pair.Value})
		}
		pod.NICs[string(result.Metric["interface"])] = nic
	}

	// network traffic receive packet
	expression.Metrics = []string{"container_network_receive_packets_total"}
	expression.QueryLabels = map[string]string{"container_label_io_kubernetes_pod_name": id}

	str = basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	str = rateExpr(durationExpr(str, int(rs.Rate)))
	resultMatrix, err = queryRange(sp, str, rs)
	if err != nil {
		return pod, err
	}

	for _, result := range resultMatrix {
		nic := pod.NICs[string(result.Metric["interface"])]
		for _, pair := range result.Values {
			nic.NICNetworkTraffic.ReceivePacketsTotal = append(nic.NICNetworkTraffic.ReceivePacketsTotal, entity.SamplePair{Timestamp: pair.Timestamp, Value: pair.Value})
		}
		pod.NICs[string(result.Metric["interface"])] = nic
	}

	// network traffic receive packet
	expression.Metrics = []string{"container_network_transmit_packets_total"}
	expression.QueryLabels = map[string]string{"container_label_io_kubernetes_pod_name": id}

	str = basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	str = rateExpr(durationExpr(str, int(rs.Rate)))
	resultMatrix, err = queryRange(sp, str, rs)
	if err != nil {
		return pod, err
	}

	for _, result := range resultMatrix {
		nic := pod.NICs[string(result.Metric["interface"])]
		for _, pair := range result.Values {
			nic.NICNetworkTraffic.TransmitPacketsTotal = append(nic.NICNetworkTraffic.TransmitPacketsTotal, entity.SamplePair{Timestamp: pair.Timestamp, Value: pair.Value})
		}
		pod.NICs[string(result.Metric["interface"])] = nic
	}

	return pod, nil
}

// GetContainer will get container
func GetContainer(sp *serviceprovider.Container, podId string, containerId string, rs RangeSetting) (entity.ContainerMetrics, error) {
	container := entity.ContainerMetrics{}

	// basic info from kube-state-metrics
	expression := Expression{}
	expression.Metrics = []string{
		"kube_pod_container_info",
		"kube_pod_container_status_restarts_total"}
	expression.QueryLabels = map[string]string{"container": containerId, "pod": podId}

	str := basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	results, err := query(sp, str)
	if err != nil {
		return container, err
	}

	for _, result := range results {
		switch result.Metric["__name__"] {

		case "kube_pod_container_info":
			container.Detail.ContainerName = containerId
			container.Detail.Pod = string(result.Metric["pod"])
			container.Detail.Image = string(result.Metric["image"])
			container.Detail.Namespace = string(result.Metric["namespace"])

		case "kube_pod_container_status_restarts_total":
			container.Detail.RestartCount = int(result.Value)
		}
	}

	// basic info from cadvisor
	expression = Expression{}
	expression.Metrics = []string{
		"container_start_time_seconds",
		"container_last_seen"}
	expression.QueryLabels = map[string]string{"container_label_io_kubernetes_container_name": containerId, "container_label_io_kubernetes_pod_name": podId}

	str = basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	results, err = query(sp, str)
	if err != nil {
		return container, err
	}

	for _, result := range results {
		switch result.Metric["__name__"] {

		case "container_last_seen":
			container.Detail.Node = string(result.Metric["node"])

		case "container_start_time_seconds":
			container.Detail.CreatedAt = int(result.Value)
		}
	}

	// status
	container.Detail.Status = ""

	expression.Metrics = []string{"kube_pod_container_status.*"}
	expression.QueryLabels = map[string]string{"container": containerId, "pod": podId}

	str = basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	str = equalExpr(str, 1)

	results, err = query(sp, str)
	if err != nil {
		return container, err
	}

containerStatusCheckingLoop:
	for _, result := range results {
		switch result.Metric["__name__"] {

		case "kube_pod_container_status_ready":
			container.Detail.Status = "ready"

		case "kube_pod_container_status_running":
			container.Detail.Status = "running"

		case "kube_pod_container_status_terminated_reason":
			container.Detail.Status = string(result.Metric["reason"])
			break containerStatusCheckingLoop

		case "kube_pod_container_status_waiting_reason":
			container.Detail.Status = string(result.Metric["reason"])
			break containerStatusCheckingLoop
		}
	}

	// command
	kc := sp.KubeCtl
	pod, err := kc.GetPod(container.Detail.Pod, container.Detail.Namespace)
	if err != nil {
		return entity.ContainerMetrics{}, err
	}

	for _, obj := range pod.Spec.Containers {
		if obj.Name == containerId {
			container.Detail.Command = obj.Command
			break
		}
	}

	if container.Detail.Status != "ready" && container.Detail.Status != "running" {
		return container, nil
	}

	// Memory resource
	expression.Metrics = []string{"container_memory_usage_bytes"}
	expression.QueryLabels = map[string]string{"container_label_io_kubernetes_container_name": containerId, "container_label_io_kubernetes_pod_name": podId}

	str = basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)

	resultMatrix, err := queryRange(sp, str, rs)
	if err != nil {
		return container, err
	}

	for _, pair := range resultMatrix[0].Values {
		container.Resource.MemoryUsageBytes = append(container.Resource.MemoryUsageBytes, entity.SamplePair{Timestamp: pair.Timestamp, Value: pair.Value})
	}

	// CPU resource
	expression.Metrics = []string{"container_cpu_usage_seconds_total"}
	expression.QueryLabels = map[string]string{"container_label_io_kubernetes_container_name": containerId, "container_label_io_kubernetes_pod_name": podId}

	str = basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	str = multiplyExpr(sumExpr(rateExpr(durationExpr(str, int(rs.Rate)))), 100)

	resultMatrix, err = queryRange(sp, str, rs)
	if err != nil {
		return container, err
	}

	for _, pair := range resultMatrix[0].Values {
		container.Resource.CPUUsagePercentage = append(container.Resource.CPUUsagePercentage, entity.SamplePair{Timestamp: pair.Timestamp, Value: pair.Value})
	}

	return container, nil
}

// GetService will get service by serviceprovider.Container
func GetService(sp *serviceprovider.Container, id string) (entity.ServiceMetrics, error) {
	service := entity.ServiceMetrics{}
	service.Labels = map[string]string{}

	// basic info
	expression := Expression{}
	expression.Metrics = []string{
		"kube_service_info",
		"kube_service_created",
		"kube_service_labels",
		"kube_service_spec_type"}
	expression.QueryLabels = map[string]string{"service": id}

	str := basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	results, err := query(sp, str)
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
		return entity.ServiceMetrics{}, fmt.Errorf("fail to get the ports setting of service %v in namespace %v: %v", service.ServiceName, service.Namespace, err)
	}
	service.Ports = object.Spec.Ports

	return service, nil
}

// GetController willl get container
func GetController(sp *serviceprovider.Container, id string) (entity.ControllerMetrics, error) {
	controller := entity.ControllerMetrics{}
	controller.Labels = map[string]string{}
	controller.Type = "deployment"

	// basic info
	expression := Expression{}
	expression.Metrics = []string{
		"kube_deployment_metadata_generation",
		"kube_deployment_created",
		"kube_deployment_labels",
		"kube_deployment_spec_replicas",
		"kube_deployment_status_replicas",
		"kube_deployment_status_replicas_available"}
	expression.QueryLabels = map[string]string{"deployment": id}

	str := basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	results, err := query(sp, str)
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

	// pods
	expression = Expression{}
	expression.Metrics = []string{"kube_pod_info"}
	expression.QueryLabels = map[string]string{"created_by_kind": "ReplicaSet", "created_by_name": id + "-[^-]+"}

	str = basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	results, err = query(sp, str)
	if err != nil {
		return controller, err
	}

	for _, result := range results {
		controller.Pods = append(controller.Pods, string(result.Metric["pod"]))
	}

	return controller, nil
}

// GetNode will get node metrics
func GetNode(sp *serviceprovider.Container, id string, rs RangeSetting) (entity.NodeMetrics, error) {
	node := entity.NodeMetrics{}
	node.Detail.Labels = map[string]string{}
	node.NICs = map[string]entity.NICMetrics{}

	// basic info
	expression := Expression{}
	expression.Metrics = []string{
		"kube_node_info",
		"kube_node_created",
		"kube_node_labels",
		"kube_node_status_capacity.*",
		"kube_node_status_allocatable.*"}
	expression.QueryLabels = map[string]string{"exported_node": id}

	str := basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	results, err := query(sp, str)
	if err != nil {
		return node, err
	}

	for _, result := range results {
		switch result.Metric["__name__"] {

		case "kube_node_info":
			node.Detail.Hostname = id
			node.Detail.KernelVersion = string(result.Metric["kernel_version"])
			node.Detail.DockerVersion = strings.Split(string(result.Metric["container_runtime_version"]), "//")[1]
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

		case "kube_node_status_allocatable_cpu_cores":
			node.Resource.AllocatableCPU = float32(result.Value)

		case "kube_node_status_allocatable_memory_bytes":
			node.Resource.AllocatableMemory = float32(result.Value)

		case "kube_node_status_allocatable_pods":
			node.Resource.AllocatablePods = float32(result.Value)

		case "kube_node_status_capacity_cpu_cores":
			node.Resource.CapacityCPU = float32(result.Value)

		case "kube_node_status_capacity_memory_bytes":
			node.Resource.CapacityMemory = float32(result.Value)

		case "kube_node_status_capacity_pods":
			node.Resource.CapacityPods = float32(result.Value)
		}
	}

	// nics
	expression = Expression{}
	expression.Metrics = []string{"node_network_interface"}
	expression.QueryLabels = map[string]string{"node": id}

	str = basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	results, err = query(sp, str)
	if err != nil {
		return node, err
	}

	for _, result := range results {
		nic := entity.NICMetrics{}
		nic.Type = string(result.Metric["type"])
		nic.PCIID = string(result.Metric["pci_id"])
		nic.IP = string(result.Metric["ip_address"])
		nic.NICNetworkTraffic = entity.NICNetworkTrafficMetrics{}
		defaultValue, err := strconv.ParseBool(string(result.Metric["default"]))
		if err != nil {
			return node, err
		}
		nic.Default = defaultValue
		dpdkValue, err := strconv.ParseBool(string(result.Metric["dpdk"]))
		if err != nil {
			return node, err
		}
		nic.DPDK = dpdkValue

		node.NICs[string(result.Metric["device"])] = nic
	}

	// status
	expression = Expression{}
	expression.Metrics = []string{"kube_node_status_condition"}
	expression.QueryLabels = map[string]string{"exported_node": id, "status": "true"}

	str = basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	str = equalExpr(str, 1)
	results, err = query(sp, str)
	if err != nil {
		return node, err
	}

	node.Detail.Status = string(results[0].Metric["condition"])

	// resource reported from kube-state-metrics
	expression = Expression{}
	expression.Metrics = []string{
		"kube_pod_container_resource_limits.*",
		"kube_pod_container_resource_requests.*"}
	expression.QueryLabels = map[string]string{"exported_node": id}
	expression.SumByLabels = []string{"__name__", "resource"}

	str = basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	str = sumByExpr(str, expression.SumByLabels)
	results, err = query(sp, str)
	if err != nil {
		return node, err
	}

	for _, result := range results {
		switch result.Metric["__name__"] {
		case "kube_pod_container_resource_requests_cpu_cores":
			node.Resource.CPURequests = float32(result.Value)

		case "kube_pod_container_resource_requests_memory_bytes":
			node.Resource.MemoryRequests = float32(result.Value)

		case "kube_pod_container_resource_limits_cpu_cores":
			node.Resource.CPULimits = float32(result.Value)

		case "kube_pod_container_resource_limits_memory_bytes":
			node.Resource.MemoryLimits = float32(result.Value)
		}
	}

	// resource resported form node-exporter
	expression = Expression{}
	expression.Metrics = []string{
		"node_memory_HugePages_Total",
		"node_memory_HugePages_Free",
		"node_memory_Hugepagesize_bytes"}
	expression.QueryLabels = map[string]string{"node": id}

	str = basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	results, err = query(sp, str)
	if err != nil {
		return node, err
	}

	for _, result := range results {
		switch result.Metric["__name__"] {
		case "node_memory_HugePages_Total":
			node.Resource.MemoryTotalHugepages = float32(result.Value)

		case "node_memory_HugePages_Free":
			node.Resource.MemoryFreeHugepages = float32(result.Value)

		case "node_memory_Hugepagesize_bytes":
			node.Resource.MemoryHugepageSize = float32(result.Value)
		}
	}

	// network traffic receive bytes
	expression.Metrics = []string{"node_network_receive_bytes_total"}
	expression.QueryLabels = map[string]string{"node": id}

	str = basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	str = rateExpr(durationExpr(str, int(rs.Rate)))
	resultMatrix, err := queryRange(sp, str, rs)
	if err != nil {
		return node, err
	}

	for _, result := range resultMatrix {
		nic := node.NICs[string(result.Metric["device"])]
		for _, pair := range result.Values {
			nic.NICNetworkTraffic.ReceiveBytesTotal = append(nic.NICNetworkTraffic.ReceiveBytesTotal, entity.SamplePair{Timestamp: pair.Timestamp, Value: pair.Value})
		}
		node.NICs[string(result.Metric["device"])] = nic
	}

	// network traffic transmit bytes
	expression.Metrics = []string{"node_network_transmit_bytes_total"}
	expression.QueryLabels = map[string]string{"node": id}

	str = basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	str = rateExpr(durationExpr(str, int(rs.Rate)))
	resultMatrix, err = queryRange(sp, str, rs)
	if err != nil {
		return node, err
	}

	for _, result := range resultMatrix {
		nic := node.NICs[string(result.Metric["device"])]
		for _, pair := range result.Values {
			nic.NICNetworkTraffic.TransmitBytesTotal = append(nic.NICNetworkTraffic.TransmitBytesTotal, entity.SamplePair{Timestamp: pair.Timestamp, Value: pair.Value})
		}
		node.NICs[string(result.Metric["device"])] = nic
	}

	// network traffic receive packets
	expression.Metrics = []string{"node_network_receive_packets_total"}
	expression.QueryLabels = map[string]string{"node": id}

	str = basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	str = rateExpr(durationExpr(str, int(rs.Rate)))
	resultMatrix, err = queryRange(sp, str, rs)
	if err != nil {
		return node, err
	}

	for _, result := range resultMatrix {
		nic := node.NICs[string(result.Metric["device"])]
		for _, pair := range result.Values {
			nic.NICNetworkTraffic.ReceivePacketsTotal = append(nic.NICNetworkTraffic.ReceivePacketsTotal, entity.SamplePair{Timestamp: pair.Timestamp, Value: pair.Value})
		}
		node.NICs[string(result.Metric["device"])] = nic
	}

	// network traffic transmit packets
	expression.Metrics = []string{"node_network_transmit_packets_total"}
	expression.QueryLabels = map[string]string{"node": id}

	str = basicExpr(expression.Metrics)
	str = queryExpr(str, expression.QueryLabels)
	str = rateExpr(durationExpr(str, int(rs.Rate)))
	resultMatrix, err = queryRange(sp, str, rs)
	if err != nil {
		return node, err
	}

	for _, result := range resultMatrix {
		nic := node.NICs[string(result.Metric["device"])]
		for _, pair := range result.Values {
			nic.NICNetworkTraffic.TransmitPacketsTotal = append(nic.NICNetworkTraffic.TransmitPacketsTotal, entity.SamplePair{Timestamp: pair.Timestamp, Value: pair.Value})
		}
		node.NICs[string(result.Metric["device"])] = nic
	}

	return node, nil
}
