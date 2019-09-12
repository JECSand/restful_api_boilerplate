/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/19/2019
*/

package server

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
	"rest_api/pkg"
	"rest_api/pkg/configuration"
)

// Server is a struct that stores the API Apps high level attributes such as the router, config, and services
type Server struct {
	Router       *mux.Router
	Config       configuration.Configuration
	UserService  root.UserService
	GroupService root.GroupService
	TodoService  root.TodoService
}

// NewServer is a function used to initialize a new Server struct
func NewServer(u root.UserService, g root.GroupService, t root.TodoService, config configuration.Configuration, client *mongo.Client) *Server {
	router := mux.NewRouter().StrictSlash(true)
	router = NewGroupRouter(g, router, config, client)
	router = NewUserRouter(u, router, g, config, client)
	router = NewTodoRouter(t, router, config, client)
	s := Server{Router: router, Config: config, UserService: u, GroupService: g, TodoService: t}
	return &s
}

// Start starts the initialized server
func (s *Server) Start() {
	if s.Config.HTTPS == "on" {
		log.Println("Listening on port 8443")
		if err := http.ListenAndServeTLS(":8443", s.Config.Cert, s.Config.Key, handlers.LoggingHandler(os.Stdout, s.Router)); err != nil {
			log.Fatal("http.ListenAndServe: ", err)
		}
	} else {
		log.Println("Listening on port 8080")
		if err := http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, s.Router)); err != nil {
			log.Fatal("http.ListenAndServe: ", err)
		}
	}
}
