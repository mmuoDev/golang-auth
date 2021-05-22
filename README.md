# golang-auth

If you are building microservices, you will surely need one of those services to handle authentication/authorization of users across other services. 
This service implements a role-based access control (RBAC) system for users. This project was built using Golang, MongoDB and Redis. 

## Functionalities
- Creation of users
- Generation of access/refresh tokens
- Uses redis to properly validate/invalidate bearer tokens 
- Uses a JSON file to define the role-based permissions for all services/resources
- Authenticate users across other services
- Logout users

## Requirements
- Golang
- MongoDB
- Redis

## Usage
Clone project and `cd` into project foler

### Starting server
``` bash
$ make run
```  

### Running Tests
``` bash
$ make test
```

## More details

### Defining role-based permissions across resources/services
- The `rbac.json` allows you define role-based permissions. For instance, to allow only `admin` access to the update functionality of a `partners`
service, you can define this:

```json
{
    "resource": "/partners",
    "methods": ["put"],
    "roles": ["admin"]
}
```
PS: If your resource url has a query parameter, use a `$` to indicate that in `rbac.json`. See the `rbac.json` file for an example

### Authenticate users across other services
- For example, to validate that a user has access to POST in the `partners` service, call the `IsAuthenticated` function in the `pkg` module. Pass the
`*http.Request` instance as a parameter. It is expected that this request includes the bearer token generated during authentication. This function returns an integer.
Each integer represents a valid http status code  except when it's 0, indicating the user is authorized. You should check for this in your service for example,

```golang
func TestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if auth := pkg.IsAuthenticated(r); auth != 0 {
			if auth == 400 {
				w.WriteHeader(http.StatusBadRequest)
			}
			//Test for other status codes
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Write([]byte("hello world!"))
	}
}
```

### Add a User
- To add a user, call the `/users` endpoint. See sample JSON request below - 
```json
{
    "firstName": "helen",
    "lastName": "ebere",
    "email": "mmuodev@gmail.com",
    "password": "password",
    "role": "admin",
    "phoneNumber": "08067170799"
}
```

### Authenticate a user
- To authenticate a user, the `phoneNumber` and `password` are required. 
```json
{
    "phoneNumber": "08067170799",
    "password": "password"
}
```
On successful authentication, amongst other things, access and refresh tokens
are generated. 
```json
{
    "id": "e27e0af8-b904-4e04-8f8c-a73db52002c7",
    "phoneNumber": "08067170799",
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjE5NGI4OGYwL...........",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjIzMDAzNjcsInJlZnJlc.........."
}
```
### Refresh a token
There is also a functionality to refresh tokens to ensure better User Experience (UX) for your users. 

# Contact
This is just a way I figured handling authentication/authorization across microservices. It can always be better. Kindly open an issue if you see ways of 
improving this. 
You can also reach out - radioactive.uche11@gmail.com


