package backend

import (
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/utils/timeutils"
	"github.com/linkernetworks/vortex/src/entity"
	"gopkg.in/mgo.v2/bson"
)

func CreateDefaultUser(mongoService *mongo.Service) error {
	session := mongoService.NewSession()
	defer session.Close()

	hashedPassword, err := HashPassword("password")
	if err != nil {
		return err
	}

	user := entity.User{
		ID: bson.NewObjectId(),
		LoginCredential: entity.LoginCredential{
			Username: "admin@vortex.com",
			Password: hashedPassword,
		},
		DisplayName: "administrator",
		Role:        "root",
		FirstName:   "administrator",
		LastName:    "administrator",
		PhoneNumber: "09521111111",
		CreatedAt:   timeutils.Now(),
	}
	q := bson.M{"loginCredential.username": user.LoginCredential.Username}

	count, err := session.Count(entity.UserCollectionName, q)
	if err != nil {
		return err
	} else if count > 0 {
		// admin user has already exists. Do not insert
		return nil
	}
	return session.Insert(entity.UserCollectionName, &user)
}
