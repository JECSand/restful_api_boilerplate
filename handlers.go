/*
Author: Connor Sanders
RESTful API Boilerplate v0.0.1
2/28/2019
*/

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}


/*
Auth Handlers
*/

/*
Test with this curl command:

curl -H "Content-Type: application/json" -d '{"username":"test", "password":"testing", "email":"test@test.com"}' http://localhost:3000/signup

*/
func Signup(w http.ResponseWriter, r *http.Request) {
	var user User
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &user); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	u := RegisterUser(user)
	if u.Email == "Taken" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(403)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusForbidden, Text: "Email Taken"}); err != nil {
			panic(err)
		}
	} else if u.Username == "Taken" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(403)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusForbidden, Text: "Username Taken"}); err != nil {
			panic(err)
		}
	} else {
		sessionToken := CreateToken(u)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Add("Auth-Token", sessionToken)
		w.WriteHeader(http.StatusCreated)
		u.Password = ""
		if err := json.NewEncoder(w).Encode(u); err != nil {
			panic(err)
		}
	}
}

/*
Test with this curl command:

curl -H "Content-Type: application/json" -d '{"username":"test", "password": "testing"}' -i http://localhost:3000/signin

*/
func Signin(w http.ResponseWriter, r *http.Request) {
	var user User
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &user); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	u := AuthenticateUser(user)
	if u.Username == "NotFound" || u.Password == "Incorrect" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(401)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusUnauthorized, Text: "Incorrect"}); err != nil {
			panic(err)
		}
	} else {
		sessionToken := CreateToken(u)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Add("Auth-Token", sessionToken)
		w.WriteHeader(http.StatusOK)
		u.Password = ""
		if err := json.NewEncoder(w).Encode(u); err != nil {
			panic(err)
		}
	}
}

func Signout(w http.ResponseWriter, r *http.Request) {

}


/*
Todos Handlers
*/


func TodoIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	todos := FindTodos()
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	todo := FindTodo(todoId)
	if todo.Uuid != "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(todo); err != nil {
			panic(err)
		}
		return
	}
	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

/*
Test with this curl command:

curl -H "Content-Type: application/json" -H "Auth-Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTE1MTQ5MDUsInJvbGUiOiJtZW1iZXIiLCJ1c2VybmFtZSI6InRlc3QyIiwidXVpZCI6ImY3MWQ2ZjNhLThlYjAtNGI2My04NTA0LWZkNDRmN2JkNGVkYyJ9.0uAKfSpL-Fxkwp-HlKGmwTxVuIOoAyLX-kpMWItvibU" -d '{"name":"New Todo2"}' http://localhost:3000/todos

*/
func TodoCreate(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &todo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	t := CreateTodo(todo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}

/*
Test with this curl command:

curl  -H "Content-Type: application/json" -X "DELETE" http://localhost:3000/todos/5fc97b70-d94b-418b-810c-e5b2d8004c9f

*/
func TodoDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	todo := DeleteTodo(todoId)
	if todo.Uuid != "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode("Todo Deleted"); err != nil {
			panic(err)
		}
		return
	}
	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}