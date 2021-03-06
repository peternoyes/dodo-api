package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		router.Methods(route.Method, "OPTIONS").Path(route.Pattern).Name(route.Name).Handler(handler)
	}

	return router
}

var routes = Routes{
	Route{
		"Build",
		"POST",
		"/build",
		Build,
	},
	Route{
		"Code",
		"GET",
		"/code/{id}",
		Code,
	},
	Route{
		"User",
		"GET",
		"/user",
		UserInfo,
	},
	Route{
		"Login",
		"GET",
		"/login",
		Login,
	},
	Route{
		"Logout",
		"GET",
		"/logout",
		Logout,
	},
	Route{
		"Callback",
		"GET",
		"/callback",
		Callback,
	},
	Route{
		"ProjectsList",
		"GET",
		"/projects",
		ProjectsList,
	},
	Route{
		"ProjectGet",
		"GET",
		"/projects/{title}",
		ProjectGet,
	},
	Route{
		"ProjectUpdate",
		"PUT",
		"/projects/{title}",
		ProjectUpdate,
	},
	Route{
		"ProjectAdd",
		"POST",
		"/projects/{title}",
		ProjectAdd,
	},
	Route{
		"ProjectDelete",
		"DELETE",
		"/projects/{title}",
		ProjectDelete,
	},
}
