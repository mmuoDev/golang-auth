package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v7"
	pkgErr "github.com/pkg/errors"
)

type TokenMetaData struct {
	AccessUUID string
	UserID     string
	Role       string
}

type RefreshTokenMeta struct {
	RefreshUUID string
	UserID      string
	Role        string
}

func CheckAuthentication(r *http.Request, client *redis.Client) (*TokenMetaData, error) {
	token, err := GetTokenMetaData(r)
	if err != nil {
		return &TokenMetaData{}, pkgErr.Wrap(err, "Token not valid")
	}
	if isValid := isTokenStoredInRedis(token, client); !isValid {
		return &TokenMetaData{}, errors.New("Token is expired")
	}
	return token, nil
}

//IsAuthenticated checks if a user is authenticated
func IsAuthenticated(r *http.Request, client *redis.Client) error {
	token, err := GetTokenMetaData(r)
	if err != nil {
		return pkgErr.Wrap(err, "Token not valid")
	}
	if isValid := isTokenStoredInRedis(token, client); !isValid {
		return errors.New("Token is expired")
	}
	return nil
}

//getToken returns token from the header
func getToken(r *http.Request) (string, error) {
	bearerToken := r.Header.Get("Authorization")
	s := strings.Split(bearerToken, " ")
	if len(s) == 2 {
		return s[1], nil
	}
	return "", errors.New("jwt - Token not found")
}

//verifyToken verifies signing mtd
func verifyToken(r *http.Request) (*jwt.Token, error) {
	ts, err := getToken(r)
	if err != nil {
		return nil, pkgErr.Wrap(err, "jwt - Token not found")
	}
	token, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, pkgErr.Wrapf(err, "Invalid token!")
	}
	return token, nil
}

//verifyRefreshToken verifies a refresh token
func verifyRefreshToken(tk string) (*jwt.Token, error) {
	token, err := jwt.Parse(tk, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})
	if err != nil {
		return nil, pkgErr.Wrap(err, "Refresh token invalid")
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, errors.New("Refresh token is invalid")
	}
	return token, nil
}

//GetRefreshTokenMetaData gets metadata off a refresh token
func GetRefreshTokenMetaData(tk string) (RefreshTokenMeta, error) {
	token, err := verifyRefreshToken(tk)
	if err != nil {
		return RefreshTokenMeta{}, pkgErr.Wrap(err, "jwt - unable to verify refresh token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		refreshUUID, ok := claims["refresh_uuid"].(string)
		if !ok {
			return RefreshTokenMeta{}, fmt.Errorf("Refresh UUID metadata not found in token")
		}
		userID, ok := claims["user_id"].(string)
		if !ok {
			return RefreshTokenMeta{}, fmt.Errorf("userID metadata not found in token")
		}
		role, ok := claims["role"].(string)
		if !ok {
			return RefreshTokenMeta{}, fmt.Errorf("Role metadata not found in token")
		}
		return RefreshTokenMeta{
			RefreshUUID: refreshUUID,
			UserID:      userID,
			Role:        role,
		}, nil
	}
	return RefreshTokenMeta{}, err
}

//isTokenValid checks if token is valid
func isTokenValid(r *http.Request) bool {
	token, err := verifyToken(r)
	if err != nil {
		return false
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return false
	}
	return true
}

//getTokenMetaData extracts metadata from a token
func GetTokenMetaData(r *http.Request) (*TokenMetaData, error) {
	token, err := verifyToken(r)
	if err != nil {
		return &TokenMetaData{}, pkgErr.Wrap(err, "jwt - unable to retrieve token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, fmt.Errorf("Access UUID metadata not found in token")
		}
		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, fmt.Errorf("userID metadata not found in token")
		}
		role, ok := claims["role"].(string)
		if !ok {
			return nil, fmt.Errorf("Role metadata not found in token")
		}
		return &TokenMetaData{
			AccessUUID: accessUUID,
			Role:       role,
			UserID:     userID,
		}, nil
	}
	return &TokenMetaData{}, err
}

//isTokenStoredInRedis checks if token still exists in redis
func isTokenStoredInRedis(td *TokenMetaData, client *redis.Client) bool {
	_, err := client.Get(td.AccessUUID).Result()
	if err != nil {
		return false
	}
	return true
}
