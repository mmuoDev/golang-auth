package app

import (
	"golang-auth/internal/db"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mmuoDev/commons/mongo"
)

//App contains handlers for the app
type App struct {
	RegisterUserHandler http.HandlerFunc
	AuthenticateHandler http.HandlerFunc
	TestHandler         http.HandlerFunc
}

//Handler returns the main handler for this application
func (a App) Handler() http.HandlerFunc {
	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/users", a.RegisterUserHandler)
	router.HandlerFunc(http.MethodPost, "/auth", a.AuthenticateHandler)
	router.HandlerFunc(http.MethodGet, "/test", a.TestHandler)

	return http.HandlerFunc(router.ServeHTTP)
}

// Options is a type for application options to modify the app
type Options func(o *Option)

// /OptionalArgs optional arguments for this application
type Option struct {
	AddUser      db.AddUserFunc
	RetrieveUser db.RetrieveUserByPhoneNumberFunc
}

//New creates a new instance of the App
func New(dbProvider mongo.DbProviderFunc, options ...Options) App {
	redisConfig := RedisInit()
	o := Option{
		AddUser:      db.AddUser(dbProvider),
		RetrieveUser: db.RetrieveUserByPhoneNumber(dbProvider),
	}

	for _, option := range options {
		option(&o)
	}

	addUser := RegisterUserHandler(o.AddUser)
	authenticate := AuthenticateHandler(o.RetrieveUser, redisConfig)
	testHandler := TestHandler(redisConfig)

	return App{
		RegisterUserHandler: addUser,
		AuthenticateHandler: authenticate,
		TestHandler:         testHandler,
	}
}
