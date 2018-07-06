package entity

type PhysicalPort struct {
	Name     string  `bson:"name" json:"name"`
	MTU      int     `bson:"MTU" json:"MTU"`
	VlanTags []int32 `bson:"vlanTags" MTC:"vlanTags"`
}

type OVSNetwork struct {
	BridgeName    string         `bson:"bridgeName" json:"bridgeName"`
	PhysicalPorts []PhysicalPort `bson:"physicalPorts" json:"physicalPorts"`
}
