package entity

type ControllerMetrics struct {
	ControllerName string            `json:"controllerName"`
	Type           string            `json:"type"`
	Namespace      string            `json:"namespace"`
	Status         string            `json:"status"`
	CreateAt       int               `json:"createAt"`
	IP             string            `json:"ip"`
	Labels         map[string]string `json:"labels"`
	RestartCount   int               `json:"restartCount"`
}
