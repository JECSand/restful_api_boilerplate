/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/19/2019
*/


package root


type Group struct {
	Id                string                   `json:"id,omitempty"`
	GroupType         string                   `json:"grouptype,omitempty"`
	Uuid              string                   `json:"uuid,omitempty"`
	Name              string                   `json:"name,omitempty"`
	LastModified      string                   `json:"last_modified,omitempty"`
	CreationDatetime  string                   `json:"creation_datetime,omitempty"`
}

type GroupService interface {
	GroupCreate(g Group) Group
	GroupFind(id string) Group
	GroupsFind() []Group
	GroupDelete(id string) Group
	GroupUpdate(g Group) Group
}
