package main

import (
	"fmt"
	"log"
	"net/http"
)

var client = DatabaseConn()

func main() {
	router := NewRouter()
	fmt.Println("Listening on port 8080....")
	log.Fatal(http.ListenAndServe(":8080", router))
}
