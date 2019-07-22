/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/19/2019
*/


package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rest_api/pkg"
)


type todoModel struct {
	Id                 primitive.ObjectID       `bson:"_id,omitempty"`
	Uuid               string                   `bson:"uuid,omitempty"`
	Name               string                   `bson:"name,omitempty"`
	Completed          string                   `bson:"completed,omitempty"`
	Due                string                   `bson:"due,omitempty"`
	Description        string                   `bson:"description,omitempty"`
	UserUuid           string                   `bson:"useruuid,omitempty"`
	GroupUuid          string                   `bson:"groupuuid,omitempty"`
	LastModified       string                   `bson:"last_modified,omitempty"`
	CreationDatetime   string                   `bson:"creation_datetime,omitempty"`
}


func newTodoModel(t root.Todo) *todoModel {
	return &todoModel{
		Uuid: t.Uuid,
		Name: t.Name,
		Completed: t.Completed,
		Due: t.Due,
		Description: t.Description,
		UserUuid: t.UserUuid,
		GroupUuid: t.GroupUuid,
		LastModified: t.LastModified,
		CreationDatetime: t.CreationDatetime,
	}
}


func(t *todoModel) toRootTodo() root.Todo {
	return root.Todo{
		Id: t.Id.Hex(),
		Uuid: t.Uuid,
		Name: t.Name,
		Completed: t.Completed,
		Due: t.Due,
		Description: t.Description,
		UserUuid: t.UserUuid,
		GroupUuid: t.GroupUuid,
		LastModified: t.LastModified,
		CreationDatetime: t.CreationDatetime,
	}
}