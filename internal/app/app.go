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
}

//Handler returns the main handler for this application
func (a App) Handler() http.HandlerFunc {
	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/users", a.RegisterUserHandler)

	return http.HandlerFunc(router.ServeHTTP)
}

// Options is a type for application options to modify the app
type Options func(o *Option)

// /OptionalArgs optional arguments for this application
type Option struct {
	AddUser db.AddUserFunc
}

//New creates a new instance of the App
func New(dbProvider mongo.DbProviderFunc, options ...Options) App {

	o := Option{
		AddUser: db.AddUser(dbProvider),
	}

	for _, option := range options {
		option(&o)
	}

	addUser := RegisterUserHandler(o.AddUser)

	return App{
		RegisterUserHandler: addUser,
	}
}
