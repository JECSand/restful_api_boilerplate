/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/19/2019
*/

package server

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"io/ioutil"
	"net/http"
	"rest_api/pkg"
	"rest_api/pkg/configuration"
)

type todoRouter struct {
	todoService root.TodoService
	config      configuration.Configuration
}

// NewTodoRouter is a function that initializes a new todoRouter struct
func NewTodoRouter(t root.TodoService, router *mux.Router, config configuration.Configuration, client *mongo.Client) *mux.Router {
	todoRouter := todoRouter{t, config}
	router.HandleFunc("/todos", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/todos", MemberTokenVerifyMiddleWare(todoRouter.TodosShow, config, client)).Methods("GET")
	router.HandleFunc("/todos/{todoId}", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/todos/{todoId}", MemberTokenVerifyMiddleWare(todoRouter.TodoShow, config, client)).Methods("GET")
	router.HandleFunc("/todos", MemberTokenVerifyMiddleWare(todoRouter.TodoCreate, config, client)).Methods("POST")
	router.HandleFunc("/todos/{todoId}", MemberTokenVerifyMiddleWare(todoRouter.TodoDelete, config, client)).Methods("DELETE")
	router.HandleFunc("/todos/{todoId}", MemberTokenVerifyMiddleWare(todoRouter.TodoModify, config, client)).Methods("PATCH")
	return router
}

// Handler function that returns all object file stored
func (tr *todoRouter) TodoModify(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	var todo root.Todo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &todo); err != nil {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	todo.Uuid = todoId
	t := tr.todoService.TodoUpdate(todo)
	if t.Uuid == "Not Found" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(404)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Todo Not Found"}); err != nil {
			panic(err)
		}
	} else {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusAccepted)
		if err := json.NewEncoder(w).Encode(t); err != nil {
			panic(err)
		}
	}
}

// Handler function that returns all object file stored
func (tr *todoRouter) TodosShow(w http.ResponseWriter, r *http.Request) {
	w = SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	decodedToken := DecodeJWT(r.Header.Get("Auth-Token"), tr.config)
	objects := tr.todoService.TodosFind(decodedToken)
	if err := json.NewEncoder(w).Encode(objects); err != nil {
		panic(err)
	}
}

// Handler function that returns a specific object file
func (tr *todoRouter) TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	decodedToken := DecodeJWT(r.Header.Get("Auth-Token"), tr.config)
	todo := tr.todoService.TodoFind(decodedToken, todoId)
	if todo.Uuid != "" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(todo); err != nil {
			panic(err)
		}
		return
	}
	// If we didn't find it, 404
	w = SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound,
		Text: "Todo Not Found"}); err != nil {
		panic(err)
	}
}

// Handler function that creates a new object and stores the file in GridFS Bucket
func (tr *todoRouter) TodoCreate(w http.ResponseWriter, r *http.Request) {
	var todo root.Todo
	decodedToken := DecodeJWT(r.Header.Get("Auth-Token"), tr.config)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &todo); err != nil {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	curid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	todo.Uuid = curid.String()
	todo.Completed = "false"
	todo.UserUuid = decodedToken[0]
	todo.GroupUuid = decodedToken[2]
	t := tr.todoService.TodoCreate(todo)
	w = SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}

// Handler function that deletes a file object from the database
func (tr *todoRouter) TodoDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	decodedToken := DecodeJWT(r.Header.Get("Auth-Token"), tr.config)
	todo := tr.todoService.TodoDelete(decodedToken, todoId)
	if todo.Uuid != "" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode("File Object Deleted"); err != nil {
			panic(err)
		}
		return
	}
	// If we didn't find it, 404
	w = SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "File Object Not Found"}); err != nil {
		panic(err)
	}
}
