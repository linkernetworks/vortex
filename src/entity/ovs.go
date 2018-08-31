package entity

// PortStatsTransmit contains information regarding the number of transmitted
// packets, bytes, etc.
type OVSPortStats struct {
	Packets uint64 `json:"packets"`
	Bytes   uint64 `json:"bytes"`
	Dropped uint64 `json:"dropped"`
	Errors  uint64 `json:"errors"`
}

type OVSPortInfo struct {
	PortID        int32        `json:"portID"`
	Name          string       `json:"name"`
	PodName       string       `json:"podName"`
	InterfaceName string       `json:"interfaceName"`
	MacAddress    string       `json:"macAddress"`
	Received      OVSPortStats `json:"received"`
	Transmitted   OVSPortStats `json:"traansmitted"`
}
