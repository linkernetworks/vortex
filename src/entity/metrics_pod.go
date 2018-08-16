package entity

// NICShortMetrics is the structure for NIC metrics
type NICShortMetrics struct {
	NICNetworkTraffic NICNetworkTrafficMetrics `json:"nicNetworkTraffic"`
}

// PodMetrics is the structure for Pod metrics
type PodMetrics struct {
	PodName      string                     `json:"podName"`
	Namespace    string                     `json:"namespace"`
	Node         string                     `json:"node"`
	UID          string                     `json:"uid"`
	Status       string                     `json:"status"`
	CreateAt     int                        `json:"createAt"`
	CreateByKind string                     `json:"createByKind"`
	CreateByName string                     `json:"createByName"`
	IP           string                     `json:"ip"`
	Labels       map[string]string          `json:"labels"`
	RestartCount int                        `json:"restartCount"`
	Containers   []string                   `json:"containers"`
	NICs         map[string]NICShortMetrics `json:"nics"`
}
