/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/19/2019
*/


package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rest_api/pkg"
	"time"
)


// GroupService is used by the app to manage all group related controllers and functionality
type GroupService struct {
	collection *mongo.Collection
	client     *mongo.Client
}


// NewGroupService is an exported function used to initialize a new GroupService struct
func NewGroupService(client *mongo.Client, dbName string, collectionName string) *GroupService {
	collection := client.Database(dbName).Collection(collectionName)
	return &GroupService {collection, client}
}


// GroupCreate is used to create a new user group
func (p *GroupService) GroupCreate(group root.Group) root.Group {
	var checkGroup = newGroupModel(root.Group{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	groupNameErr := p.collection.FindOne(ctx, bson.M{"name": group.Name}).Decode(&checkGroup)
	if groupNameErr == nil {
		return root.Group{Name: "Taken"}
	}
	currentTime := time.Now().UTC()
	group.LastModified = currentTime.String()
	group.CreationDatetime = currentTime.String()
	groupModel := newGroupModel(group)
	_, err2 := p.collection.InsertOne(ctx, groupModel)
	if err2 != nil {
		fmt.Println("group doc creation error err: ", err2)
	}
	return groupModel.toRootGroup()
}


// GroupsFind is used to find all group docs
func (p *GroupService) GroupsFind() []root.Group {
	var groups []root.Group
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := p.collection.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var group = newGroupModel(root.Group{})
		cursor.Decode(&group)
		groups = append(groups, group.toRootGroup())
	}
	return groups
}


// GroupFind is used to find a specific group doc
func (p *GroupService) GroupFind(id string) root.Group {
	var group = newGroupModel(root.Group{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := p.collection.FindOne(ctx, bson.M{"uuid": id}).Decode(&group)
	if err != nil {
		return root.Group{}
	}
	return group.toRootGroup()
}


// GroupDelete is used to delete a group doc
func (p *GroupService) GroupDelete(id string) root.Group {
	var group = newGroupModel(root.Group{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := p.collection.FindOneAndDelete(ctx, bson.M{"uuid": id}).Decode(&group)
	if err != nil {
		return root.Group{}
	}
	return group.toRootGroup()
}


// GroupUpdate is used to update an existing group
func (p *GroupService) GroupUpdate(group root.Group) root.Group {
	var curGroup = newGroupModel(root.Group{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	groupErr := p.collection.FindOne(ctx, bson.M{"uuid": group.Uuid}).Decode(&curGroup)
	if groupErr != nil {
		return root.Group{Uuid: "Not Found"}
	}
	filter := bson.D{{"uuid", group.Uuid}}
	currentTime := time.Now().UTC()
	update := bson.D{{"$set",
		bson.D{
			{"name", group.Name},
			{"last_modified", currentTime.String()},
		},
	}}
	_, err := p.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("update err: ", err)
		panic(err)
	}
	return group
}


// GroupDocInsert is used to insert a group doc directly into mongodb for testing purposes
func (p *GroupService) GroupDocInsert(group root.Group) root.Group {
	var insertGroup = newGroupModel(group)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := p.collection.InsertOne(ctx, insertGroup)
	if err != nil {
		fmt.Println("group doc insertion error: ", err)
	}
	return insertGroup.toRootGroup()
}