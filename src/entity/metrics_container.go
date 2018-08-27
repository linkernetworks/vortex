package entity

// ContainerResourceMetrics is the structure for Container Resource Metrics
type ContainerResourceMetrics struct {
	CPUUsagePercentage []SamplePair `json:"cpuUsagePercentage"`
	MemoryUsageBytes   []SamplePair `json:"memoryUsageBytes"`
}

// ContainerDetailMetrics is the structure  for Container Detail Metrics
type ContainerDetailMetrics struct {
	ContainerName string   `json:"containerName"`
	CreatedAt     int      `json:"createAt"`
	Status        string   `json:"status"`
	RestartCount  int      `json:"restartCount"`
	Pod           string   `json:"pod"`
	Namespace     string   `json:"namespace"`
	Node          string   `json:"node"`
	Image         string   `json:"image"`
	Command       []string `json:"command"`
}

// ContainerMetrics is the structure for Container Metrics
type ContainerMetrics struct {
	Detail   ContainerDetailMetrics   `json:"detail"`
	Resource ContainerResourceMetrics `json:"resource"`
}
