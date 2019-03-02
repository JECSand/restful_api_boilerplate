/*
Author: Connor Sanders
RESTful API Boilerplate v0.0.1
2/28/2019
*/

package main

import (
	"fmt"
	"log"
	"net/http"
)

var client = DatabaseConn()

func main() {
	router := NewRouter()
	fmt.Println("Listening on port 3000....")
	log.Fatal(http.ListenAndServe(":3000", router))
}
