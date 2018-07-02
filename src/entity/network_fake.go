package entity

type FakeNetwork struct {
	FakeParameter string `bson:"bridgeName" json:"bridgeName"`
	IWantFail     bool   `bson:"iWantFail" json:"iWantFail"`
}
