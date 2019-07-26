/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/25/2019
*/

package root

// Blacklist is a root struct that is used to store the json encoded data for/from a mongodb blacklist doc.
type Blacklist struct {
	Id               string `json:"id,omitempty"`
	AuthToken        string `json:"auth_token,omitempty"`
	LastModified     string `json:"last_modified,omitempty"`
	CreationDatetime string `json:"creation_datetime,omitempty"`
}
