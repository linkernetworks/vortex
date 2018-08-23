package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// UserCollectionName's const
const (
	UserCollectionName string = "users"
)

// RegistryBasicAuthCredential is the structure for a user login credential
type RegistryBasicAuthCredential struct {
	Username string `bson:"username" json:"username" validate:"required"`
	Password string `bson:"password" json:"password" validate:"required"`
}

// LoginCredential is the structure for a user login credential
type LoginCredential struct {
	Username string `bson:"username" json:"username" validate:"required,email"`
	Password string `bson:"password" json:"password" validate:"required"`
}

// User is the structure for user info
type User struct {
	ID              bson.ObjectId   `bson:"_id,omitempty" json:"id" validate:"-"`
	LoginCredential LoginCredential `bson:"loginCredential" json:"loginCredential" validate:"required"`
	DisplayName     string          `bson:"displayName" json:"displayName" validate:"required"`
	Role            string          `bson:"role" json:"role" validate:"required,eq=root|eq=user|eq=guest"`
	FirstName       string          `bson:"firstname" json:"firstName" validate:"required"`
	LastName        string          `bson:"lastName" json:"lastName" validate:"required"`
	PhoneNumber     string          `bson:"phoneNumber" json:"phoneNumber" validate:"required,numeric"`
	CreatedAt       *time.Time      `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"-"`
}

// GetCollection - get model mongo collection name.
func (u User) GetCollection() string {
	return UserCollectionName
}
