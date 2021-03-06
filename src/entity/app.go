package entity

// Deployment is the structure for deployment info
type Application struct {
	Deployment Deployment `bson:"deployment" json:"deployment" validate:"required"`
	Service    Service    `bson:"service" json:"service" validate:"required"`
}
