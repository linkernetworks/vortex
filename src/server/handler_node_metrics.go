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

//vortex-dev
func getNodeMetricsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	id := req.PathParameter("id")

	results, _ := queryFromPrometheus(sp, `{__name__=~"kube_node_info|kube_node_created|node_network_interface|kube_node_labels|kube_node_status_allocatable_cpu_cores|kube_node_status_allocatable_memory_bytes|kube_node_status_capacity_memory_bytes|",node=~"`+id+`"}`)

	node := entity.NodeMetrics{}
	nics := &node.Info.NICs
	labels := &node.Info.Labels

	for _, result := range results {
		switch result.Metric["__name__"] {

		case "kube_node_info":
			node.Info.Hostname = id
			node.Info.KernelVersion = string(result.Metric["kernel_version"])
			node.Info.OS = string(result.Metric["os_image"])
			node.Info.KubernetesVersion = string(result.Metric["kubelet_version"])

		case "kube_node_created":
			node.Info.CreatedAt = int(result.Value)

		case "node_network_interface":
			nic := entity.NICMetrics{}
			nic.Name = string(result.Metric["device"])
			nic.Default = string(result.Metric["default"])
			nic.Type = string(result.Metric["type"])
			nic.IP = string(result.Metric["ip_address"])
			*nics = append(*nics, nic)

		case "kube_node_labels":
			label := entity.NodeLabelMetrics{}
			for key, value := range result.Metric {
				if strings.HasPrefix(string(key), "label_") {
					label.Key = strings.TrimPrefix(string(key), "label_")
					label.Value = string(value)
					*labels = append(*labels, label)
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
	// results, _ := queryFromPrometheus(sp, `{__name__=~"kube_node_info",node=~"`+id+`"}`)
	// node.Info.Hostname = id
	// node.Info.KernelVersion = string(results[0].Metric["kernel_version"])
	// node.Info.OS = string(results[0].Metric["os_image"])
	// node.Info.KubernetesVersion = string(results[0].Metric["kubelet_version"])

	// results, _ = queryFromPrometheus(sp, `{__name__=~"kube_node_created",node=~"`+id+`"}`)
	// node.Info.CreatedAt = int(results[0].Value)

	// results, _ = queryFromPrometheus(sp, `{__name__=~"node_network_interface",node=~"`+id+`"}`)
	// NICs := []entity.NetworkInterfaceMetrics{}
	// for _, result := range results {
	// 	nic := entity.NetworkInterfaceMetrics{}
	// 	nic.Name = string(result.Metric["device"])
	// 	nic.Default = string(result.Metric["default"])
	// 	nic.Type = string(result.Metric["type"])
	// 	nic.IP = string(result.Metric["ip_address"])
	// 	NICs = append(NICs, nic)
	// }
	// node.Info.NICs = NICs

	//node.Info.NIC = string(results[0].Value)

	resp.WriteEntity(node)

	// sum(kube_node_status_capacity_cpu_cores)-sum(kube_pod_container_resource_requests{resource="cpu",node="vortex-dev"})
	// sum (kube_pod_container_resource_limits{resource="cpu",node="vortex-dev"})
	// sum (kube_pod_container_resource_requests{resource="cpu",node="vortex-dev"})

	// sum (kube_pod_container_resource_limits{resource="memory",node="vortex-dev"})
	// sum (kube_pod_container_resource_requests{resource="memory",node="vortex-dev"})

}
