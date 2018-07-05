package entity

type NFSStorage struct {
	IP   string `bson:"ip" json:"ip"`
	PATH string `bson:"path" json:"path"`
}
