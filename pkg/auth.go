package pkg

import (
	"golang-auth/internal"
	"golang-auth/internal/middleware"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-redis/redis/v7"
	"github.com/mmuoDev/commons/httputils"
)

//IsAuthenticated checks if user is authenticated and authorized to use a resource
func IsAuthenticated(r *http.Request) int {
	client := RedisInit()
	token, err := middleware.CheckAuthentication(r, client)
	if err != nil {
		log.Printf("Token is invalid for http request=%v", r)
		return http.StatusUnauthorized
	}
	role := token.Role
	url := r.URL.Path
	method := r.Method

	var rbac []internal.RBAC

	httputils.FileToStruct(filepath.Join("rbac.json"), &rbac)
	roles, methods, isURL := getRBAC(rbac, url)

	if !isURL {
		log.Printf("URL is undefined!, url=%s", url)
		return http.StatusBadRequest
	}

	isMethod := checkHTTPMethod(methods, method)
	if !isMethod {
		log.Printf("Method is not permitted to access resource=%s, method=%s", url, method)
		return http.StatusUnauthorized
	}

	isRole := checkRole(roles, role)
	if !isRole {
		log.Printf("Role is not permitted to access resource, role=%s, resource=%s", role, url)
		return http.StatusUnauthorized
	}
	return 0
}

//getRBAC validates URL/resource and returns the associated http methods and roles
func getRBAC(rbac []internal.RBAC, url string) ([]string, []string, bool) {
	sURL := strings.Split(url, "/")
	for _, v := range rbac {
		r := v.Resource
		mtds := v.Methods
		roles := v.Roles

		sRes := strings.Split(r, "/")
		if len(sRes) == len(sURL) {
			i := 0
			for _, r := range sRes {
				if r == sURL[i] || r == "$" {
					if i == len(sRes)-1 {
						return roles, mtds, true
					}
					i = i + 1
					continue
				} else {
					break
				}
			}
		}
	}
	return []string{}, []string{}, false
}

//checkRole validates if a role is authorized to access a resource
func checkRole(roles []string, role string) bool {
	for _, r := range roles {
		if strings.ToLower(r) == strings.ToLower(role) {
			return true
		}
	}
	return false
}

//checkHTTPMethod validates if a http method is authorized to access a resource
func checkHTTPMethod(methods []string, mtd string) bool {
	for _, m := range methods {
		if strings.ToLower(m) == strings.ToLower(mtd) {
			return true
		}
	}
	return false
}

func RedisInit() *redis.Client {
	var client *redis.Client
	client = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_DSN"),
	})
	return client
}
