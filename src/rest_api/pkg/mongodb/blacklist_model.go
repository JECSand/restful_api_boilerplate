/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/25/2019
*/


package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rest_api/pkg"
)


type blacklistModel struct {
	Id                primitive.ObjectID       `bson:"_id,omitempty"`
	AuthToken         string                   `bson:"auth_token,omitempty"`
	LastModified      string                   `bson:"last_modified,omitempty"`
	CreationDatetime  string                   `bson:"creation_datetime,omitempty"`
}


func newBlacklistModel(bl root.Blacklist) *blacklistModel {
	return &blacklistModel{
		AuthToken: bl.AuthToken,
		LastModified: bl.LastModified,
		CreationDatetime: bl.CreationDatetime,
	}
}


func(bl *blacklistModel) toRootBlacklist() root.Blacklist {
	return root.Blacklist{
		Id: bl.Id.Hex(),
		AuthToken: bl.AuthToken,
		LastModified: bl.LastModified,
		CreationDatetime: bl.CreationDatetime,
	}
}
