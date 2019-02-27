package main

import (
	"context"
	"fmt"
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
	err := collection.FindOne(ctx, Todo{Guid: id}).Decode(&todo)
	if err != nil {
		panic(err)
	}
	return todo
}


func RepoCreateTodo(t Todo) Todo {
	// currentId needs to be a UUID
	currentId += 1
	t.Guid = currentId
	collection := client.Database("testing").Collection("todos")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection.InsertOne(ctx, t)
	return t
}


func RepoDestroyTodo(id int) error {
	var todos Todos
	for i, t := range todos {
		if t.Guid == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Todo with id of %d to delete", id)
}
