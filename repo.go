package main

import (
	"context"
	"fmt"
	"time"
)

var currentId int
var todos Todos


// Give us some seed data
//func init() {
	//RepoCreateTodo(Todo{Name: "Write presentation"})
	//RepoCreateTodo(Todo{Name: "Host meetup"})
//}

func RepoFindTodo(id int) Todo {
	for _, t := range todos {
		if t.Guid == id {
			return t
		}
	}
	// return empty Todo if not found
	return Todo{}
}

//this is bad, I don't think it passes race condtions
func RepoCreateTodo(t Todo) Todo {
	// currentId needs to be a UUID
	currentId += 1
	t.Guid = currentId
	fmt.Println(client)
	collection := client.Database("testing").Collection("todos")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, t)
	fmt.Println(result)
	return t
}

func RepoDestroyTodo(id int) error {
	for i, t := range todos {
		if t.Guid == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Todo with id of %d to delete", id)
}
