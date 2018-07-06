package entity

type DPDKPhysicalPort struct {
	Name     string  `bson:"name" json:"name"`
	MTU      int     `bson:"MTU" json:"MTU"`
	PCIID    string  `bson:"pciID" json:"pciID"`
	VlanTags []int32 `bson:"vlanTags" json:"vlanTags"`
}

type OVSDPDKNetwork struct {
	BridgeName        string             `bson:"bridgeName" json:"bridgeName"`
	DPDKPhysicalPorts []DPDKPhysicalPort `bson:"physicalPorts" json:"physicalPorts"`
}
