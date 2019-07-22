/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/19/2019
*/


package root


type Todo struct {
	Id                string                   `json:"id,omitempty"`
	Uuid              string                   `json:"uuid,omitempty"`
	Name              string                   `json:"name,omitempty"`
	Completed         string                   `json:"completed,omitempty"`
	Due               string                   `json:"due,omitempty"`
	Description       string                   `json:"description,omitempty"`
	UserUuid          string                   `json:"useruuid,omitempty"`
	GroupUuid         string                   `json:"groupuuid,omitempty"`
	LastModified      string                   `json:"last_modified,omitempty"`
	CreationDatetime  string                   `json:"creation_datetime,omitempty"`
}

type TodoService interface {
	TodoCreate(t Todo) Todo
	TodoFind([]string, string) Todo
	TodosFind([]string) []Todo
	TodoDelete([]string, string) Todo
	TodoUpdate(t Todo) Todo
}
