package entity

type NetworkInterface struct {
	Name    string
	Default bool
	IP      string
}

type Node struct {
	Name          string
	CPU           float32
	Memory        float32
	Pods          int
	InterfaceList []NetworkInterface
	Labels        []map[string]string
}
