package app

import (
	"golang-auth/internal/db"
	"golang-auth/internal/middleware"
	"golang-auth/internal/workflow"
	"golang-auth/pkg"
	"log"
	"net/http"
	"path/filepath"

	"golang-auth/internal"

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
func AuthenticateHandler(retrieveUser db.RetrieveUserByPhoneNumberFunc, client *redis.Client) http.HandlerFunc {
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

//LogoutHandler returns  http request to logout a user
func LogoutHandler(client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		td, err := middleware.GetTokenMetaData(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		logout := workflow.Logout(client)
		if err := logout(td.AccessUUID); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

//RefreshTokenHandler returns http requst to refresh tokens
func RefreshTokenHandler(client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rt pkg.RefreshTokenRequest
		httputils.JSONToDTO(&rt, w, r)

		refresh := workflow.RefreshToken(client)
		d, err := refresh(rt)
		if err != nil {
			httputils.ServeError(err, w)
			return
		}
		httputils.ServeJSON(d, w)
	}
}

func IsAuthenticated(r *http.Request, w http.ResponseWriter, client *redis.Client) string {
	token, err := middleware.CheckAuthentication(r, client)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}
	role := token.Role
	url := r.URL.Path
	method := r.Method

	var rbac []internal.RBAC
	httputils.FileToStruct(filepath.Join("rbac.json"), &rbac)

	for _, v := range rbac {
		r := v.Resource
		mtds := v.Methods
		roles := v.Roles

		
		log.Fatal(k, v)
	}
	return role + url + method
}

func TestHandler(client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := IsAuthenticated(r, w, client)
		w.Write([]byte(s))
		// if err := middleware.IsAuthenticated(r, client); err != nil {
		// 	log.Fatal( r.URL.Path)

		// 	//Testing RBAC
		// 	//TODO: Put in the pkg folder
		// 	var rbac internal.RBAC
		// 	httputils.JSONToDTO(&rbac, w, r)

		// 	w.Write([]byte("unauthenticated"))
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	return
		// }
		// w.Write([]byte("Welcome!"))
	}
}
