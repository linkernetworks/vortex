package entity

// physical NIC port
type PhysicalPort struct {
	Name     string  `bson:"name" json:"name"`
	MTU      int     `bson:"MTU" json:"MTU"`
	VlanTags []int32 `bson:"vlanTags" MTC:"vlanTags"`
}

// physical NIC port dpdk enabled
type DPDKPhysicalPort struct {
	Name     string  `bson:"name" json:"name"`
	MTU      int     `bson:"MTU" json:"MTU"`
	PCIID    string  `bson:"pciID" json:"pciID"`
	VlanTags []int32 `bson:"vlanTags" json:"vlanTags"`
}

// kernel space datapath
type OVSNetwork struct {
	BridgeName    string         `bson:"bridgeName" json:"bridgeName"`
	PhysicalPorts []PhysicalPort `bson:"physicalPorts" json:"physicalPorts"`
}

// userspace space datapath
type OVSUserspaceNetwork struct {
	BridgeName string `bson:"bridgeName" json:"bridgeName"`
	// exclusive fields
	PhysicalPorts     []PhysicalPort     `bson:"physicalPorts" json:"physicalPorts"`
	DPDKPhysicalPorts []DPDKPhysicalPort `bson:"dpdkPhysicalPorts" json:"dpdkPhysicalPorts"`
}
