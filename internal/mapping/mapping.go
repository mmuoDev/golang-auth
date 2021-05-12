package mapping

import (
	"golang-auth/internal"
	"golang-auth/pkg"

	"github.com/mmuoDev/commons/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

//generateHashPassword generates password hash from a string
func generateHashPassword(password string) (string, error) {
	bb, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", errors.Wrap(err, "Mapping- unable to hash password")
	}
	return string(bb), nil
}

//ToDBUser maps user to db user
func ToDBUser(u pkg.User) (internal.User, error) {
	password, err := generateHashPassword(u.Password)
	if err != nil {
		return internal.User{}, errors.Wrap(err, "Mapping - unable to hash password")
	}
	return internal.User{
		ID:          uuid.GenV4(),
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Email:       u.Email,
		Password:    password,
		Role:        u.Role,
		PhoneNumber: u.PhoneNumber,
		IsVerified:  false,
	}, nil
}