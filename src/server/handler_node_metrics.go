package server

import (
	"fmt"
	"strings"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/vortex/src/entity"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/web"
	"github.com/prometheus/common/model"
)

func listNodeMetricsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	result, err := queryFromPrometheus(sp, "sum by (node)(kube_node_info)")

	if result == nil {
		response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("%v: %v", result, err))
	}

	nodeList := []model.LabelValue{}

	for _, node := range result {
		nodeList = append(nodeList, node.Metric["node"])
	}

	logger.Infof("fetching all nodes. found %d nodes", len(nodeList))
	resp.WriteEntity(nodeList)
}

func getNodeMetricsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	id := req.PathParameter("id")
	node := entity.NodeMetrics{}
	node.Info.Labels = map[string]string{}
	node.NICs = map[string]entity.NICMetrics{}

	results, _ := queryFromPrometheus(sp, `{__name__=~"kube_node_info|kube_node_created|node_network_interface|kube_node_labels|kube_node_status_allocatable_cpu_cores|kube_node_status_allocatable_memory_bytes|kube_node_status_capacity_cpu_cores|kube_node_status_capacity_memory_bytes|",node=~"`+id+`"}`)

	for _, result := range results {
		switch result.Metric["__name__"] {

		case "kube_node_info":
			node.Info.Hostname = id
			node.Info.KernelVersion = string(result.Metric["kernel_version"])
			node.Info.OS = string(result.Metric["os_image"])
			node.Info.KubernetesVersion = string(result.Metric["kubelet_version"])

		case "kube_node_created":
			node.Info.CreatedAt = int(result.Value)

		case "kube_node_labels":
			for key, value := range result.Metric {
				if strings.HasPrefix(string(key), "label_") {
					node.Info.Labels[strings.TrimPrefix(string(key), "label_")] = string(value)
				}
			}

		case "kube_node_status_allocatable_cpu_cores":
			node.Resource.AllocatableCPU = float32(result.Value)

		case "kube_node_status_allocatable_memory_bytes":
			node.Resource.AllocatableMemory = float32(result.Value)

		case "kube_node_status_capacity_cpu_cores":
			node.Resource.CapacityCPU = float32(result.Value)

		case "kube_node_status_capacity_memory_bytes":
			node.Resource.CapacityMemory = float32(result.Value)
		}
	}

	results, _ = queryFromPrometheus(sp, `sum by(__name__, resource) ({__name__=~"kube_pod_container_resource_limits|kube_pod_container_resource_requests",node=~"`+id+`"})`)

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

	results, _ = queryFromPrometheus(sp, `{__name__=~"node_network_interface|node_network_receive_bytes_total|node_network_transmit_bytes_total|node_network_receive_packets_total|node_network_transmit_packets_total",node=~"`+id+`"}`)

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

	resp.WriteEntity(node)

}
