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


type Server struct {
	router             *mux.Router
	config             configuration.Configuration
}


func NewServer(u root.UserService, g root.GroupService, t root.TodoService, config configuration.Configuration, client *mongo.Client) *Server {
	router := mux.NewRouter().StrictSlash(true)
	router = NewGroupRouter(g, router, config, client)
	router = NewUserRouter(u, router, g, config, client)
	router = NewTodoRouter(t, router, config, client)
	s := Server{router: router, config: config}
	return &s
}


func (s *Server) Start() {
	if s.config.HTTPS == "on" {
		log.Println("Listening on port 8443")
		if err := http.ListenAndServeTLS(":8443", s.config.Cert, s.config.Key, handlers.LoggingHandler(os.Stdout, s.router)); err != nil {
			log.Fatal("http.ListenAndServe: ", err)
		}
	} else {
		log.Println("Listening on port 8080")
		if err := http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, s.router)); err != nil {
			log.Fatal("http.ListenAndServe: ", err)
		}
	}
}