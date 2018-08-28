package entity

// ControllerMetrics is the structure for Controller Metrics
type ControllerMetrics struct {
	ControllerName string            `json:"controllerName"`
	Type           string            `json:"type"`
	Namespace      string            `json:"namespace"`
	Strategy       string            `json:"strategy"`
	CreateAt       int               `json:"createAt"`
	DesiredPod     int               `json:"desiredPod"`
	CurrentPod     int               `json:"currentPod"`
	AvailablePod   int               `json:"availablePod"`
	Pods           []string          `json:"pods"`
	Labels         map[string]string `json:"labels"`
}
