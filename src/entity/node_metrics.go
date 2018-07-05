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
	AllocatableCPU              float32 `json:"allocatableCPU"`
	AllocatableMemory           float32 `json:"allocatableMemory"`
	AllocatablePods             float32 `json:"allocatablePods"`
	AllocatableEphemeralStorage float32 `json:"allocatableEphemeralStorage"`
	CapacityCPU                 float32 `json:"capacityCPU"`
	CapacityMemory              float32 `json:"capacityMemory"`
	CapacityPods                float32 `json:"capacityPods"`
	CapacityEphemeralStorage    float32 `json:"capacityEphemeralStorage"`
	CPURequests                 float32 `json:"cpuRequests"`
	CPULimits                   float32 `json:"cpuLimits"`
	MemoryRequests              float32 `json:"memoryRequests"`
	MemoryLimits                float32 `json:"memoryLimits"`
}

type NodeResourceShortMetrics struct {
	CPURequests    float32 `json:"cpuRequests"`
	CPULimits      float32 `json:"cpuLimits"`
	MemoryRequests float32 `json:"memoryRequests"`
	MemoryLimits   float32 `json:"memoryLimits"`
}

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

type NodeInfoMetrics struct {
	NodeName string                   `json:"nodeName"`
	Labels   map[string]string        `json:"labels"`
	Status   string                   `json:"status"`
	Resource NodeResourceShortMetrics `json:"resource"`
}

type NodeListMetrics struct {
	Node map[string]NodeInfoMetrics `json:"node"`
}

type NodeMetrics struct {
	Detail   NodeDetailMetrics     `json:"detail"`
	Resource NodeResourceMetrics   `json:"resource"`
	NICs     map[string]NICMetrics `json:"nics"`
}
