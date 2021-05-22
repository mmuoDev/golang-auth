package app

import (
	"fmt"
	"golang-auth/internal/db"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mmuoDev/commons/mongo"
)

const (
	aID = "accountId"
)

//App contains handlers for the app
type App struct {
	RegisterUserHandler http.HandlerFunc
	AuthenticateHandler http.HandlerFunc
	TestHandler         http.HandlerFunc
	LogoutHandler       http.HandlerFunc
	RefreshTokenHandler http.HandlerFunc
}

//Handler returns the main handler for this application
func (a App) Handler() http.HandlerFunc {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, fmt.Sprintf("/test/:%s/refresh", aID), a.TestHandler)

	router.HandlerFunc(http.MethodPost, "/users", a.RegisterUserHandler)
	router.HandlerFunc(http.MethodPost, "/auth", a.AuthenticateHandler)
	router.HandlerFunc(http.MethodPost, "/logout", a.LogoutHandler)
	router.HandlerFunc(http.MethodPost, "/token/refresh", a.RefreshTokenHandler)

	return http.HandlerFunc(router.ServeHTTP)
}

// Options is a type for the app options
type Options func(o *OptionalArgs)

// /OptionalArgs defines optional arguments for this app
type OptionalArgs struct {
	AddUser      db.AddUserFunc
	RetrieveUser db.RetrieveUserByPhoneNumberFunc
}

//New creates a new instance of the App
func New(dbProvider mongo.DbProviderFunc, options ...Options) App {
	redisConfig := RedisInit()
	o := OptionalArgs{
		AddUser:      db.AddUser(dbProvider),
		RetrieveUser: db.RetrieveUserByPhoneNumber(dbProvider),
	}

	for _, option := range options {
		option(&o)
	}

	addUser := RegisterUserHandler(o.AddUser)
	authenticate := AuthenticateHandler(o.RetrieveUser, redisConfig)
	testHandler := TestHandler(redisConfig)
	logoutHandler := LogoutHandler(redisConfig)
	refreshTokenHandler := RefreshTokenHandler(redisConfig)

	return App{
		RegisterUserHandler: addUser,
		AuthenticateHandler: authenticate,
		TestHandler:         testHandler,
		LogoutHandler:       logoutHandler,
		RefreshTokenHandler: refreshTokenHandler,
	}
}
