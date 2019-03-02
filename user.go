/*
Author: Connor Sanders
RESTful API Boilerplate v0.0.1
2/28/2019
*/

package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type User struct {
	ID            primitive.ObjectID       `json:"_id,omitempty" bson:"_id,omitempty"`
	Uuid          string                   `json:"uuid,omitempty" bson:"uuid,omitempty"`
	Username      string   				   `json:"username,omitempty" bson:"username,omitempty"`
	Password      string                   `json:"password,omitempty" bson:"password,omitempty"`
	Email         string                   `json:"email,omitempty" bson:"email,omitempty"`
	Role          string                   `json:"role,omitempty" bson:"role,omitempty"`
}

type Users []User
