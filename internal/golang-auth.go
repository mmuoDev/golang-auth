package internal

import (
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

//TokenDetails defines access and refresh tokens
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	ATExpires    int64
	RTExpires    int64
}

type HTTPMethod struct {
	Method string `json:"method"`
}

type Role struct {
	Role string `json:"role"`
}

//RBAC represents data for role-based access control
type RBAC struct {
	Resource string       `json:"resource"`
	Methods  []HTTPMethod `json:"methods"`
	Roles    []Role       `json:"roles"`
}
