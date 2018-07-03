package entity

type NICMetrics struct {
	Name    string
	Default string
	Type    string
	IP      string
}

type NodeLabelMetrics struct {
	Key   string
	Value string
}

type NodeInfoMetrics struct {
	Hostname          string
	KernelVersion     string
	CreatedAt         int
	OS                string
	KubernetesVersion string
	Labels            []NodeLabelMetrics
	NICs              []NICMetrics
}

type NodeResourceMetrics struct {
	AllocatableCPU    float32
	AllocatableMemory float32
	CapacityCPU       float32
	CapacityMemory    float32
}

type NICNetworkTrafficMetrics struct {
	TransmitBytes int
	ReveiveBytes  int
	PacketsCount  int
}

type NodeMetrics struct {
	Info              NodeInfoMetrics
	Resource          NodeResourceMetrics
	NICNetworkTraffic []NICNetworkTrafficMetrics
}
