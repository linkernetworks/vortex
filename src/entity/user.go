package entity

import (
	"time"

	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)

// UserCollectionName's const
const (
	UserCollectionName string = "users"
)

// User is the structure for user info
type User struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id" validate:"-"`
	UUID        uuid.UUID     `bson:"uuid" json:"uuid" validate:"required,uuid4"`
	JWT         string        `bson:"jwt" json:"jwt" validate:"-"`
	UserName    string        `bson:"username" json:"username" validate:"required"`
	Email       string        `bson:"email" json:"email" validate:"required,email"`
	Password    string        `bson:"password" json:"password" validate:"required"`
	FirstName   string        `bson:"firstname" json:"firstName" validate:"required"`
	LastName    string        `bson:"lastName" json:"lastName" validate:"required"`
	PhoneNumber string        `bson:"phoneNumber" json:"phoneNumber" validate:"required,numeric"`
	CreatedAt   *time.Time    `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"-"`
}

// GetCollection - get model mongo collection name.
func (u User) GetCollection() string {
	return UserCollectionName
}
