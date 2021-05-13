package app

import (
	"golang-auth/internal/db"
	"golang-auth/internal/workflow"
	"golang-auth/pkg"
	"net/http"

	"github.com/go-redis/redis/v7"
	"github.com/mmuoDev/commons/httputils"
)

//RegisterUserHandler returns a http request to register a  user
func RegisterUserHandler(addUser db.AddUserFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user pkg.User
		httputils.JSONToDTO(&user, w, r)
		//TODO: Validation

		add := workflow.AddUser(addUser)
		if err := add(user); err != nil {
			httputils.ServeError(err, w)
			return 
		}
		w.WriteHeader(http.StatusCreated)
	}
}

//AuthenticateHandler returns a http request to authenticate a user
func AuthenticateHandler (retrieveUser db.RetrieveUserByPhoneNumberFunc, client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var a pkg.Login
		httputils.JSONToDTO(&a, w, r)
		//TODO: Validation

		auth := workflow.Authenticate(retrieveUser, client)
		u, err := auth(a)
		if err != nil {
			httputils.ServeError(err, w)
			return 
		}
		httputils.ServeJSON(u, w)
	}
}