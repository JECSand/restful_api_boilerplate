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
	"rest_api/pkg/configuration"
	"time"
)

// TodoService is used by the app to manage all todos related controllers and functionality
type TodoService struct {
	collection *mongo.Collection
	config     configuration.Configuration
	client     *mongo.Client
}

// NewTodoService is an exported function used to initialize a new TodoService struct
func NewTodoService(client *mongo.Client, dbName string, collectionName string, config configuration.Configuration) *TodoService {
	collection := client.Database(dbName).Collection(collectionName)
	return &TodoService{collection, config, client}
}

// TodosFind is used to Find all todos docs
func (p *TodoService) TodosFind(decodedToken []string) []root.Todo {
	var todos []root.Todo
	var queryFilter = newTodoModel(root.Todo{})
	if decodedToken[1] != "master_admin" {
		queryFilter.GroupUuid = decodedToken[2]
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := p.collection.Find(ctx, queryFilter)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var todo = newTodoModel(root.Todo{})
		cursor.Decode(&todo)
		todos = append(todos, todo.toRootTodo())
	}
	return todos
}

// TodoFind is used to Find a specific todos doc
func (p *TodoService) TodoFind(decodedToken []string, id string) root.Todo {
	var todo = newTodoModel(root.Todo{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := p.collection.FindOne(ctx, bson.M{"uuid": id}).Decode(&todo)
	if decodedToken[1] != "master_admin" && todo.GroupUuid != decodedToken[2] {
		return root.Todo{GroupUuid: "DoesNotMatch"}
	}
	if err != nil {
		return root.Todo{}
	}
	return todo.toRootTodo()
}

// TodoCreate is used to create a new todos doc
func (p *TodoService) TodoCreate(todo root.Todo) root.Todo {
	currentTime := time.Now().UTC()
	todo.LastModified = currentTime.String()
	todo.CreationDatetime = currentTime.String()
	todoModel := newTodoModel(todo)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := p.collection.InsertOne(ctx, todoModel)
	if err != nil {
		fmt.Println("todo doc creation error err: ", err)
	}
	return todoModel.toRootTodo()
}

// TodoDelete is used to delete a specified todos doc
func (p *TodoService) TodoDelete(decodedToken []string, id string) root.Todo {
	var todo = newTodoModel(root.Todo{})
	filter := bson.M{"uuid": id}
	if decodedToken[1] != "master_admin" {
		groupUuid := decodedToken[2]
		userUuid := decodedToken[0]
		filter = bson.M{"uuid": id, "useruuid": userUuid, "groupuuid": groupUuid}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := p.collection.FindOneAndDelete(ctx, filter).Decode(&todo)
	if err != nil {
		return root.Todo{}
	}
	return todo.toRootTodo()
}

// TodoUpdate is used to update an existing todos doc
func (p *TodoService) TodoUpdate(todo root.Todo) root.Todo {
	var curTodo = newTodoModel(root.Todo{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	todoErr := p.collection.FindOne(ctx, bson.M{"uuid": todo.Uuid}).Decode(&curTodo)
	if todoErr != nil {
		return root.Todo{Uuid: "Not Found"}
	}
	if len(todo.Name) == 0 {
		todo.Name = curTodo.Name
	}
	if len(todo.Due) == 0 {
		todo.Due = curTodo.Due
	}
	if len(todo.Description) == 0 {
		todo.Description = curTodo.Description
	}
	if len(todo.Completed) == 0 {
		todo.Completed = curTodo.Completed
	}
	if len(todo.UserUuid) == 0 {
		todo.UserUuid = curTodo.UserUuid
	}
	filter := bson.D{{"uuid", todo.Uuid}}
	currentTime := time.Now().UTC()
	update := bson.D{{"$set",
		bson.D{
			{"name", todo.Name},
			{"due", todo.Due},
			{"description", todo.Description},
			{"completed", todo.Completed},
			{"useruuid", todo.UserUuid},
			{"last_modified", currentTime.String()},
		},
	}}
	_, err := p.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("update err: ", err)
		panic(err)
	}
	return todo
}

// TodoDocInsert is used to insert a todos doc directly into mongodb for testing purposes
func (p *TodoService) TodoDocInsert(todo root.Todo) root.Todo {
	var insertTodo = newTodoModel(todo)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := p.collection.InsertOne(ctx, insertTodo)
	if err != nil {
		fmt.Println("todo doc insertion error: ", err)
	}
	return insertTodo.toRootTodo()
}
