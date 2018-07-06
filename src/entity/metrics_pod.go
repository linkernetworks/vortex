package entity

type PodDetailMetrics struct {
	PodName      string            `json:"podName"`
	Namespace    string            `json:"namespace"`
	Node         string            `json:"node"`
	Status       string            `json:"status"`
	CreateAt     int               `json:"createAt"`
	Labels       map[string]string `json:"labels"`
	RestartCount int               `json:"restartCount"`
}

type PodMetrics struct {
	Detail     PodDetailMetrics `json:"detail"`
	Containers string           `json:"container"`
	Events     string           `json:"events"`
}
