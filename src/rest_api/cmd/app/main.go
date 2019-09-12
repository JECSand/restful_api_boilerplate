/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/19/2019
*/

package main

func main() {
	a := App{}
	a.Initialize("production")
	a.Run()
}
