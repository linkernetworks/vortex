package entity

type NICNetworkTrafficMetrics struct {
	ReceiveBytesTotal    int `json:"receiveBytesTotal"`
	TransmitBytesTotal   int `json:"transmitBytesTotal"`
	ReceivePacketsTotal  int `json:"receivePacketsTotal"`
	TransmitPacketsTotal int `json:"transmitPacketsTotal"`
}

type NICMetrics struct {
	Default           string                   `json:"default"`
	Type              string                   `json:"type"`
	IP                string                   `json:"ip"`
	NICNetworkTraffic NICNetworkTrafficMetrics `json:"nicNetworkTraffic"`
}

type NodeResourceMetrics struct {
	AllocatableCPU    float32 `json:"allocatableCPU"`
	AllocatableMemory float32 `json:"allocatableMemory"`
	CapacityCPU       float32 `json:"capacityCPU"`
	CapacityMemory    float32 `json:"capacityMemory"`
	CPURequests       float32 `json:"cPURequests"`
	CPULimits         float32 `json:"cPULimits"`
	MemoryRequests    float32 `json:"memoryRequests"`
	MemoryLimits      float32 `json:"memoryLimits"`
}

type NodeInfoMetrics struct {
	Hostname          string            `json:"hostname"`
	KernelVersion     string            `json:"kernelVersion"`
	CreatedAt         int               `json:"createAt"`
	OS                string            `json:"os"`
	KubernetesVersion string            `json:"kubernetesVersion"`
	Labels            map[string]string `json:"labels"`
}

type NodeMetrics struct {
	Info     NodeInfoMetrics       `json:"info"`
	Resource NodeResourceMetrics   `json:"resource"`
	NICs     map[string]NICMetrics `json:nics""`
}
