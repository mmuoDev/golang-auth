package workflow

import (
	"golang-auth/internal"
	"golang-auth/internal/db"
	"golang-auth/internal/mapping"
	"golang-auth/internal/middleware"
	"golang-auth/pkg"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v7"
	"github.com/mmuoDev/commons/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

//AddUserFunc adds a user
type AddUserFunc func(r pkg.User) error

//AuthenticateFunc authenticates a user
type AuthenticateFunc func(r pkg.Login) (pkg.Auth, error)

//LogoutFunc provides functionality to logout an authenticated user
type LogoutFunc func (accessUUID string) error

//RefreshTokenFunc generates a new refresh token
type RefreshTokenFunc func(rt pkg.RefreshTokenRequest) (pkg.RefreshToken, error)

//RefreshToken returns functionality to generate another access/refresh token using a refresh token
func RefreshToken (client *redis.Client) RefreshTokenFunc {
	return func(rt pkg.RefreshTokenRequest) (pkg.RefreshToken, error) {
		tk, err := middleware.GetRefreshTokenMetaData(rt.RefreshToken)
		if err != nil {
			return pkg.RefreshToken{}, errors.Wrap(err, "Workflow - Invalid refresh token")
		}
		//delete previous refresh token
		d, err := deleteJWT(client, tk.RefreshUUID)
		if err != nil || d == 0 { 
			return pkg.RefreshToken{}, errors.Wrap(err, "Workflow - Expired refresh token")
		}
		//create new access and refresh tokens
		td, err := generateJWT(tk.UserID, tk.Role)
		if err != nil {
			return pkg.RefreshToken{}, errors.Wrap(err, "Workflow - error generating access and refresh tokens")
		}
		//save token to redis
		if err := saveJWTMetaData(client, td, tk.UserID); err != nil {
			return pkg.RefreshToken{}, errors.Wrap(err, "Workflow - error saving token to redis")
		}
		return mapping.ToRefreshToken(td), nil
	}
}

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
func Authenticate(retrieveUser db.RetrieveUserByPhoneNumberFunc, client *redis.Client) AuthenticateFunc {
	return func(r pkg.Login) (pkg.Auth, error) {
		retrieve, err := retrieveUser(r.PhoneNumber)
		if err != nil {
			return pkg.Auth{}, errors.Wrapf(err, "Workflow - No user found for phone=%s", r.PhoneNumber)
		}
		if vp := validatePassword(retrieve.Password, r.Password); !vp {
			return pkg.Auth{}, errors.New("Workflow - Incorrect auth credentials")
		}
		userID := retrieve.ID.Val()
		td, err := generateJWT(userID, retrieve.Role)
		if err != nil {
			return pkg.Auth{}, errors.Wrapf(err, "Workflow - Unable to generate tokens")
		}
		if err := saveJWTMetaData(client, td, userID); err != nil {
			return pkg.Auth{}, errors.Wrap(err, "Workflow - unable to save tokens in redis")
		}
		return mapping.ToAuth(retrieve, td), nil
	}
}

//Logout logouts an authenticated user
func Logout(client *redis.Client) LogoutFunc {
	return func(accessUUID string) error {
		d, err := deleteJWT(client, accessUUID) 
		if err != nil || d == 0 {
			return errors.Wrap(err, "Workflow - unable to delete id. User may be unauthorized!")
		}
		return nil
	}
}
//validatePassword validates password for a user
func validatePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//generateJWT generates a token
func generateJWT(id, role string) (*internal.TokenDetails, error) {
	td := &internal.TokenDetails{}
	td.ATExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUUID = uuid.GenV4().Val()

	td.RTExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = uuid.GenV4().Val()

	//Access token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = id
	atClaims["exp"] = td.ATExpires
	atClaims["role"] = role
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	t, err := at.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
	if err != nil {
		return nil, errors.Wrap(err, "Workflow - unable to generate access token")
	}
	td.AccessToken = t

	//Refresh token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = id
	rtClaims["exp"] = td.RTExpires
	rtClaims["role"] = role
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
	if err != nil {
		return nil, errors.Wrap(err, "Workflow - unable to generate refresh token")
	}

	return td, nil
}

//saveJWTMetaData saves JWT's meta data to redis
func saveJWTMetaData(client *redis.Client, td *internal.TokenDetails, userID string) error {
	at := time.Unix(td.ATExpires, 0) //convert Unix to UTC(to Time object)
	rt := time.Unix(td.RTExpires, 0)
	now := time.Now()

	if err := client.Set(td.AccessUUID, userID, at.Sub(now)).Err(); err != nil {
		return errors.Wrap(err, "Workflow - unable to save accessUUID to redis!")
	}
	if err := client.Set(td.RefreshUUID, userID, rt.Sub(now)).Err(); err != nil {
		return errors.Wrap(err, "Workflow - unable to save refreshUUID to redis!")
	}
	return nil
}

//deleteJWT deletes a token from redis. Deleting from redis invalidates the token
func deleteJWT(client *redis.Client, uuid string) (int64, error) {
	d, err := client.Del(uuid).Result()
	if err != nil {
		return 0, errors.Wrapf(err, "Workflow - unable to delete uuid from redis, uuid=%s", uuid)
	}
	return d, nil
}


