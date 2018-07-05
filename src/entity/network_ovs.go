package entity

type PhysicalPort struct {
	Name     string `bson:"name" json:"name"`
	MTU      int    `bson:"MTU" json:"MTU"`
	VlanTags []int  `bson:"vlanTag" json:"vlanTag"`
}

type OVSNetwork struct {
	BridgeName    string         `bson:"bridgeName" json:"bridgeName"`
	PhysicalPorts []PhysicalPort `bson:"physicalPorts" json:"physicalPorts"`
}
