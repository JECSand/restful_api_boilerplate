/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/19/2019
*/

package root

// Group is a root struct that is used to store the json encoded data for/from a mongodb group doc.
type Group struct {
	Id               string `json:"id,omitempty"`
	GroupType        string `json:"grouptype,omitempty"`
	Uuid             string `json:"uuid,omitempty"`
	Name             string `json:"name,omitempty"`
	LastModified     string `json:"last_modified,omitempty"`
	CreationDatetime string `json:"creation_datetime,omitempty"`
}

// GroupService is an interface used to manage the relevant group doc controllers
type GroupService interface {
	GroupCreate(g Group) Group
	GroupFind(id string) Group
	GroupsFind() []Group
	GroupDelete(id string) Group
	GroupUpdate(g Group) Group
	GroupDocInsert(g Group) Group
}
