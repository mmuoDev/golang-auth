package workflow

import (
	"golang-auth/internal"
	"golang-auth/internal/db"
	"golang-auth/internal/mapping"
	"golang-auth/pkg"

	"github.com/pkg/errors"
)

//AddUserFunc adds a user
type AddUserFunc func(r pkg.User) error

//AuthenticateFunc authenticates a user
type AuthenticateFunc func(r pkg.Auth) (internal.Auth, error)

//AddUser adds a user 
func AddUser(addUser db.AddUserFunc) AddUserFunc {
	return func(r pkg.User) error {
		u, err := mapping.ToDBUser(r)

		if err != nil {
			return errors.Wrap(err, "Workflow - unable to map internal user to db")
		}
		if err := addUser(u); err != nil {
			return errors.Wrap(err, "Workflow - error adding new user")
		}
		return nil
	}
}

//Authenticate authenticates a user
func Authenticate(retrieveUser db.RetrieveUserByPhoneNumberFunc) AuthenticateFunc {
	return func(r pkg.Auth) (internal.Auth, error) {

	}
}
