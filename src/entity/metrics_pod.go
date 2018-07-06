package entity

type PodDetailMetrics struct {
	PodName      string            `json:"podName"`
	Namespace    string            `json:"namespace"`
	Node         string            `json:"node"`
	Status       string            `json:"status"`
	CreateAt     int               `json:"createAt"`
	CreateByKind string            `json:"createByKind"`
	CreateByName string            `json:"createByName"`
	IP           string            `json:"ip"`
	Labels       map[string]string `json:"labels"`
	RestartCount int               `json:"restartCount"`
}

type PodMetrics struct {
	Detail     PodDetailMetrics `json:"detail"`
	Containers []string         `json:"containers"`
	Events     string           `json:"events"`
}
