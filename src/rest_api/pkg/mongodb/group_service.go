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


type GroupService struct {
	collection *mongo.Collection
	client     *mongo.Client
}


func NewGroupService(client *mongo.Client, dbName string, collectionName string) *GroupService {
	collection := client.Database(dbName).Collection(collectionName)
	return &GroupService {collection, client}
}


// Create a new user group
func (p *GroupService) GroupCreate(group root.Group) root.Group {
	var checkGroup = newGroupModel(root.Group{})
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
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


// Function to find groups
func (p *GroupService) GroupsFind() []root.Group {
	var groups []root.Group
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
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


// Find a Group
func (p *GroupService) GroupFind(id string) root.Group {
	var group = newGroupModel(root.Group{})
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := p.collection.FindOne(ctx, bson.M{"uuid": id}).Decode(&group)
	if err != nil {
		return root.Group{}
	}
	return group.toRootGroup()
}


// Delete a Group
func (p *GroupService) GroupDelete(id string) root.Group {
	var group = newGroupModel(root.Group{})
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := p.collection.FindOneAndDelete(ctx, bson.M{"uuid": id}).Decode(&group)
	if err != nil {
		return root.Group{}
	}
	return group.toRootGroup()
}


// Update an existing group
func (p *GroupService) GroupUpdate(group root.Group) root.Group {
	var curGroup = newGroupModel(root.Group{})
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	groupErr := p.collection.FindOne(ctx, bson.M{"uuid": group.Uuid}).Decode(&curGroup)
	if groupErr != nil {
		return root.Group{Uuid: "Not Found"}
	}  else {
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
	return root.Group{}
}