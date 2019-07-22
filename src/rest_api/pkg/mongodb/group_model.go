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


type groupModel struct {
	Id                primitive.ObjectID       `bson:"_id,omitempty"`
	GroupType         string                   `bson:"grouptype,omitempty"`
	Uuid              string                   `bson:"uuid,omitempty"`
	Name              string                   `bson:"name,omitempty"`
	LastModified      string                   `bson:"last_modified,omitempty"`
	CreationDatetime  string                   `bson:"creation_datetime,omitempty"`
}


func newGroupModel(g root.Group) *groupModel {
	return &groupModel{
		GroupType: g.GroupType,
		Uuid: g.Uuid,
		Name: g.Name,
		LastModified: g.LastModified,
		CreationDatetime: g.CreationDatetime,
	}
}


func(g *groupModel) toRootGroup() root.Group {
	return root.Group{
		Id: g.Id.Hex(),
		GroupType: g.GroupType,
		Uuid: g.Uuid,
		Name: g.Name,
		LastModified: g.LastModified,
		CreationDatetime: g.CreationDatetime,
	}
}