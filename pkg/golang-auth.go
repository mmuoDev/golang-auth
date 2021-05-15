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

//Login represents data needed to authenticate a user and generate a token
type Login struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

//Auth represents data after successful authentication
type Auth struct {
	ID           string `json:"id"`
	PhoneNumber  string `json:"phoneNumber"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

//RefreshTokenRequest represents request body to refresh a token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

//RefreshToken represents response for generating refresh token
type RefreshToken struct {
	AccessToken  string
	RefreshToken string
}
