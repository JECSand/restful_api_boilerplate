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

type userModel struct {
	Id               primitive.ObjectID `bson:"_id,omitempty"`
	Uuid             string             `bson:"uuid,omitempty"`
	Username         string             `bson:"username,omitempty"`
	Password         string             `bson:"password,omitempty"`
	Email            string             `bson:"email,omitempty"`
	Role             string             `bson:"role,omitempty"`
	GroupUuid        string             `bson:"groupuuid,omitempty"`
	LastModified     string             `bson:"last_modified,omitempty"`
	CreationDatetime string             `bson:"creation_datetime,omitempty"`
}

func newUserModel(u root.User) *userModel {
	return &userModel{
		Uuid:             u.Uuid,
		Username:         u.Username,
		Password:         u.Password,
		Email:            u.Email,
		Role:             u.Role,
		GroupUuid:        u.GroupUuid,
		LastModified:     u.LastModified,
		CreationDatetime: u.CreationDatetime,
	}
}

func (u *userModel) toRootUser() root.User {
	return root.User{
		Id:               u.Id.Hex(),
		Uuid:             u.Uuid,
		Username:         u.Username,
		Password:         u.Password,
		Email:            u.Email,
		Role:             u.Role,
		GroupUuid:        u.GroupUuid,
		LastModified:     u.LastModified,
		CreationDatetime: u.CreationDatetime,
	}
}
