/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/19/2019
*/

package server

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}
