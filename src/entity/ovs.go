package entity

// PortStatsReceive contains information regarding the number of received
// packets, bytes, etc.
type PortStatsReceive struct {
	Packets uint64 `json:"packets"`
	Bytes   uint64 `json:"bytes"`
	Dropped uint64 `json:"dropped"`
	Errors  uint64 `json:"errors"`
	Frame   uint64 `json:"-"`
	Over    uint64 `json:"-"`
	CRC     uint64 `json:"-"`
}

// PortStatsTransmit contains information regarding the number of transmitted
// packets, bytes, etc.
type PortStatsTransmit struct {
	Packets    uint64 `json:"packets"`
	Bytes      uint64 `json:"bytes"`
	Dropped    uint64 `json:"dropped"`
	Errors     uint64 `json:"errors"`
	Collisions uint64 `json:"collisions"`
}
type OVSPortStat struct {
	PortID      uint32            `json:'portID"`
	Received    PortStatsReceive  `json:"received"`
	Transmitted PortStatsTransmit `json:"traansmitted"`
}
