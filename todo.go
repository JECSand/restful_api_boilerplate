/*
Author: Connor Sanders
RESTful API Boilerplate v0.0.1
2/28/2019
*/

package main

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Todo struct {
	ID        primitive.ObjectID       `json:"_id,omitempty" bson:"_id,omitempty"`
	Uuid      string                   `json:"uuid,omitempty" bson:"uuid,omitempty"`
	Name      string   				   `json:"name,omitempty" bson:"name,omitempty"`
	Completed bool                     `json:"completed,omitempty" bson:"completed,omitempty"`
	Due       time.Time                `json:"due,omitempty" bson:"due,omitempty"`
}

type Todos []Todo
