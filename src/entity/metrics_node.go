package entity

import "github.com/prometheus/common/model"

// SamplePair pairs a SampleValue with a Timestamp.
type SamplePair struct {
	Timestamp model.Time        `json:"timestamp"`
	Value     model.SampleValue `json:"value"`
}

// NICNetworkTrafficMetrics is the structure for NIC metwork traffic metrics
type NICNetworkTrafficMetrics struct {
	ReceiveBytesTotal    []SamplePair `json:"receiveBytesTotal"`
	TransmitBytesTotal   []SamplePair `json:"transmitBytesTotal"`
	ReceivePacketsTotal  []SamplePair `json:"receivePacketsTotal"`
	TransmitPacketsTotal []SamplePair `json:"transmitPacketsTotal"`
}

// NICMetrics is the structure for NIC metrics
type NICMetrics struct {
	Default           bool                     `json:"default"`
	DPDK              bool                     `json:"dpdk"`
	Type              string                   `json:"type"`
	IP                string                   `json:"ip"`
	PCIID             string                   `json:"pciID"`
	NICNetworkTraffic NICNetworkTrafficMetrics `json:"nicNetworkTraffic"`
}

// NICOverviewMetrics is the structure for NIC overview metrics
type NICOverviewMetrics struct {
	Name    string `json:"name"`
	Default bool   `json:"default"`
	DPDK    bool   `json:"dpdk"`
	Type    string `json:"type"`
	PCIID   string `json:"pciID"`
}

// NodeResourceMetrics is the structure for node resource metrics
type NodeResourceMetrics struct {
	CPURequests       float32 `json:"cpuRequests"`
	CPULimits         float32 `json:"cpuLimits"`
	MemoryRequests    float32 `json:"memoryRequests"`
	MemoryLimits      float32 `json:"memoryLimits"`
	AllocatableCPU    float32 `json:"allocatableCPU"`
	AllocatableMemory float32 `json:"allocatableMemory"`
	AllocatablePods   float32 `json:"allocatablePods"`
	CapacityCPU       float32 `json:"capacityCPU"`
	CapacityMemory    float32 `json:"capacityMemory"`
	CapacityPods      float32 `json:"capacityPods"`
}

// NodeDetailMetrics is the structure for node detail metrics
type NodeDetailMetrics struct {
	Hostname          string            `json:"hostname"`
	CreatedAt         int               `json:"createAt"`
	Status            string            `json:"status"`
	OS                string            `json:"os"`
	KernelVersion     string            `json:"kernelVersion"`
	KubeproxyVersion  string            `json:"kubeproxyVersion"`
	KubernetesVersion string            `json:"kubernetesVersion"`
	Labels            map[string]string `json:"labels"`
}

// NodeNICsMetrics is the structure for node NICs metrics
type NodeNICsMetrics struct {
	NICs []NICOverviewMetrics `json:"nics"`
}

// NodeMetrics is the structure for node metrics
type NodeMetrics struct {
	Detail   NodeDetailMetrics     `json:"detail"`
	Resource NodeResourceMetrics   `json:"resource"`
	NICs     map[string]NICMetrics `json:"nics"`
}
