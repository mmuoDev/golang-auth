package internal

import (
	"time"

	"github.com/mmuoDev/commons/uuid"
)

//User represents an internal user
type User struct {
	ID          uuid.V4 `bson:"id"`
	FirstName   string  `bson:"firstName"`
	LastName    string  `bson:"lastName"`
	Email       string  `bson:"email"`
	Password    string  `bson:"password"`
	Role        string  `bson:"role"`
	PhoneNumber string  `bson:"phoneNumber"`
	IsVerified  bool    `bson:"isVerified"`
}

//Auth represents response after a successful authentication
type Auth struct {
	ID           string    `json:"id"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Role         string    `json:"role"`
	Token        string    `json:"token"`
	TokenExpires time.Time `json:"token_expires_on"`
}
