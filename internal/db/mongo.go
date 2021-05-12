package db

import (
	"golang-auth/internal"

	"github.com/mmuoDev/commons/mongo"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	usersCollection = "users"
)

//AddUserFunc returns functionality to add a user
type AddUserFunc func(internal.User) error

//RetrieveUserByPhoneNumberFunc returns functionality to retrieve a user by phone number
type RetrieveUserByPhoneNumberFunc func(phoneNumber string) (internal.User, error)

func AddUser(dbProvider mongo.DbProviderFunc) AddUserFunc {
	return func(u internal.User) error {
		col := mongo.NewCollection(dbProvider, usersCollection)
		_, err := col.Insert(u)
		if err != nil {
			return errors.Wrap(err, "db - failure inserting a user")
		}
		return nil
	}
}

//RetrieveUserByPhoneNumber retrieves user by phone number
func RetrieveUserByPhoneNumber(dbProvider mongo.DbProviderFunc) RetrieveUserByPhoneNumberFunc {
	return func(phoneNumber string) (internal.User, error) {
		col := mongo.NewCollection(dbProvider, usersCollection)
		filter := bson.D{{"phoneNumber", phoneNumber}}
		var user internal.User

		if err := col.FindOne(filter, &user); err != nil {
			return internal.User{}, errors.Wrapf(err, "db - user not found")
		}
		return user, nil
	}
}
