package pkg

//User represents a user
type User struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	PhoneNumber string `json:"phoneNumber"`
	IsVerified  bool   `json:"isVerified"`
}

//Auth represents data needed to authenticate a user and generate a token
type Auth struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}
