package server

import (
	"strings"

	"github.com/linkernetworks/vortex/src/entity"
	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/web"
)

func listNodeMetricsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	nodeList := entity.NodeListMetrics{}
	nodeList.Node = map[string]entity.NodeInfoMetrics{}

	// kube_node_info, kube_node_labels
	results, err := queryFromPrometheus(sp, `{__name__=~"kube_node_info|kube_node_labels"}`)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	for _, result := range results {
		switch result.Metric["__name__"] {

		case "kube_node_info":
			nodeInfo := entity.NodeInfoMetrics{}
			nodeInfo.NodeName = string(result.Metric["node"])
			nodeList.Node[string(result.Metric["node"])] = nodeInfo

		case "kube_node_labels":
			nodeInfo := nodeList.Node[string(result.Metric["node"])]
			nodeInfo.Labels = map[string]string{}
			for key, value := range result.Metric {
				if strings.HasPrefix(string(key), "label_") {
					nodeInfo.Labels[strings.TrimPrefix(string(key), "label_")] = string(value)
				}
			}
			nodeList.Node[string(result.Metric["node"])] = nodeInfo
		}
	}

	// kube_node_status_condition
	results, err = queryFromPrometheus(sp, `{__name__=~"kube_node_status_condition",status="true"}==1`)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	for _, result := range results {
		nodeInfo := nodeList.Node[string(result.Metric["node"])]
		nodeInfo.Status = string(result.Metric["condition"])
		nodeList.Node[string(result.Metric["node"])] = nodeInfo
	}

	// kube_pod_container_resource_limits, kube_pod_container_resource_requests
	results, err = queryFromPrometheus(sp, `sum by(__name__, resource,node) ({__name__=~"kube_pod_container_resource_limits|kube_pod_container_resource_requests"})`)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	for _, result := range results {
		nodeInfo := nodeList.Node[string(result.Metric["node"])]
		switch result.Metric["__name__"] {
		case "kube_pod_container_resource_requests":
			switch result.Metric["resource"] {
			case "cpu":
				nodeInfo.Resource.CPURequests = float32(result.Value)
			case "memory":
				nodeInfo.Resource.MemoryRequests = float32(result.Value)
			}
		case "kube_pod_container_resource_limits":
			switch result.Metric["resource"] {
			case "cpu":
				nodeInfo.Resource.CPULimits = float32(result.Value)
			case "memory":
				nodeInfo.Resource.MemoryLimits = float32(result.Value)
			}
		}
		nodeList.Node[string(result.Metric["node"])] = nodeInfo
	}

	resp.WriteEntity(nodeList)
}

func getNodeMetricsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	id := req.PathParameter("id")
	node := entity.NodeMetrics{}
	node.Detail.Labels = map[string]string{}
	node.NICs = map[string]entity.NICMetrics{}

	// kube_node_info, kube_node_created, node_network_interface, kube_node_labels, kube_node_status_capacity, kube_node_status_allocatable
	results, err := queryFromPrometheus(sp, `{__name__=~"kube_node_info|kube_node_created|node_network_interface|kube_node_labels|kube_node_status_capacity|kube_node_status_allocatable",node=~"`+id+`"}`)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
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
	results, err = queryFromPrometheus(sp, `	{__name__=~"kube_node_status_condition",node=~"`+id+`",status="true"}==1`)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	node.Detail.Status = string(results[0].Metric["condition"])

	// kube_pod_container_resource_limits, kube_pod_container_resource_requests
	results, err = queryFromPrometheus(sp, `sum by(__name__, resource) ({__name__=~"kube_pod_container_resource_limits|kube_pod_container_resource_requests",node=~"`+id+`"})`)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
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
	results, err = queryFromPrometheus(sp, `{__name__=~"node_network_interface|node_network_receive_bytes_total|node_network_transmit_bytes_total|node_network_receive_packets_total|node_network_transmit_packets_total",node=~"`+id+`"}`)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
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

	resp.WriteEntity(node)
}
