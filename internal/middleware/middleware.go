package middleware

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

//IsAuthenticated checks if a user is authenticated 
func IsAuthenticated(r *http.Request, client *redis.Client) error {
	token, err := getTokenMetaData(r)
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
	return "", errors.New("Middleware - Token not found")
}

//verifyToken verifies signing mtd
func verifyToken(r *http.Request) (*jwt.Token, error) {
	ts, err := getToken(r)
	if err != nil {
		return nil, pkgErr.Wrap(err, "Middleware - Token not found")
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
func getTokenMetaData(r *http.Request) (*TokenMetaData, error) {
	token, err := verifyToken(r)
	if err != nil {
		return &TokenMetaData{}, pkgErr.Wrap(err, "MIddleware - unable to retrieve token")
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
