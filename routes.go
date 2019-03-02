/*
Author: Connor Sanders
RESTful API Boilerplate v0.0.1
2/28/2019
*/

package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Signup",
		"POST",
		"/signup",
		Signup,
	},
	Route{
		"Signin",
		"POST",
		"/signin",
		Signin,
	},
	Route{
		"Signout",
		"DELETE",
		"/signout",
		Signout,
	},
	Route{
		"TodoIndex",
		"GET",
		"/todos",
		MemberTokenVerifyMiddleWare(TodoIndex),
	},
	Route{
		"TodoCreate",
		"POST",
		"/todos",
		MemberTokenVerifyMiddleWare(TodoCreate),
	},
	Route{
		"TodoShow",
		"GET",
		"/todos/{todoId}",
		MemberTokenVerifyMiddleWare(TodoShow),
	},
	Route{
		"TodoDelete",
		"DELETE",
		"/todos/{todoId}",
		MemberTokenVerifyMiddleWare(TodoDelete),
	},
}
