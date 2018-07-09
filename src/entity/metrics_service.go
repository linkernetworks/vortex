package entity

import (
	"k8s.io/api/core/v1"
)

type ServiceMetrics struct {
	ServiceName string            `json:"serviceName"`
	Namespace   string            `json:"namespace"`
	Type        string            `json:"type"`
	CreateAt    int               `json:"createAt"`
	ClusterIP   string            `json:"clusterIP"`
	Ports       []v1.ServicePort  `json: "ports"`
	Labels      map[string]string `json:"labels"`
}
