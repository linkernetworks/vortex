package entity

type ContainerResourceMetrics struct {
	CPUUsagePercentage float32 `json:"cpuUsagePercentage"`
	MemoryUsageBytes   float32 `json:"memoryUsageBytes"`
}

type ContainerStatusMetrics struct {
	Status           string `json:"status"`
	WaitingReason    string `json:"waitingReason"`
	TerminatedReason string `json:"terminatedReason"`
	RestartTime      int    `json:"restartTime"`
}

type ContainerDetailMetrics struct {
	ContainerName string `json:"containerName"`
	CreatedAt     int    `json:"createAt"`
	Pod           string `json:"pod"`
	Node          string `json:"node"`
	Image         string `json:"image"`
	Command       string `json:"command"`
	vNIC          string `json:"vNic"`
}
type ContainerMetrics struct {
	Detail            ContainerDetailMetrics   `json:"detail"`
	Status            ContainerStatusMetrics   `json:"status"`
	Resource          ContainerResourceMetrics `json:"resource"`
	NICNetworkTraffic NICNetworkTrafficMetrics `json:"nicNetworkTraffic"`
}
