package entity

type ServiceMetrics struct {
	ServiceName string            `json:"serviceName"`
	Namespace   string            `json:"namespace"`
	Type        string            `json:"type"`
	CreateAt    int               `json:"createAt"`
	ClusterIP   string            `json:"clusterIP"`
	Labels      map[string]string `json:"labels"`
}
