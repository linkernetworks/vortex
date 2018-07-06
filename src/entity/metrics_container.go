package entity

type ContainerMetrics struct {
	ContainerName string            `json:"containerName"`
	Namespace     string            `json:"namespace"`
	Node          string            `json:"node"`
	Status        string            `json:"status"`
	CreateAt      int               `json:"createAt"`
	CreateByKind  string            `json:"createByKind"`
	CreateByName  string            `json:"createByName"`
	IP            string            `json:"ip"`
	Labels        map[string]string `json:"labels"`
	RestartCount  int               `json:"restartCount"`
}
