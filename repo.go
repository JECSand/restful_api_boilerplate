package main

import (
	"context"
	//"fmt"
	"time"
	"go.mongodb.org/mongo-driver/bson"
)

var currentId int


func init() {
	RepoCreateTodo(Todo{Name: "Write presentation"})
	RepoCreateTodo(Todo{Name: "Host meetup"})
}


func RepoFindTodos() []Todo {
	var todos []Todo
	collection := client.Database("testing").Collection("todos")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var todo Todo
		cursor.Decode(&todo)
		todos = append(todos, todo)
	}
	return todos
}


func RepoFindTodo(id int) Todo {
	var todo Todo
	collection := client.Database("testing").Collection("todos")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Todo{Uuid: id}).Decode(&todo)
	if err != nil {
		panic(err)
	}
	return todo
}


func RepoCreateTodo(t Todo) Todo {
	// currentId needs to be a UUID
	currentId += 1
	t.Uuid = currentId
	collection := client.Database("testing").Collection("todos")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	collection.InsertOne(ctx, t)
	return t
}


func RepoDeleteTodo(id int) Todo {
	var todo Todo
	collection := client.Database("testing").Collection("todos")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOneAndDelete(ctx, Todo{Uuid: id}).Decode(&todo)
	if err != nil {
		panic(err)
	}
	return todo
}
