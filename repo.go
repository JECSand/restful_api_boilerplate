/*
Author: Connor Sanders
RESTful API Boilerplate v0.0.1
2/28/2019
*/

package main

import (
	"context"
	"fmt"
	"time"
	"github.com/gofrs/uuid"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var MySigningKey = []byte("secret")


func init() {
	CreateTodo(Todo{Name: "Write presentation"})
	CreateTodo(Todo{Name: "Host meetup"})
}

/*
Auth Utilities
*/

func CreateToken(user User) string {
	/* Create the token */
	token := jwt.New(jwt.SigningMethodHS256)

	/* Create a map to store our claims*/
	claims := token.Claims.(jwt.MapClaims)
	claims["role"] = user.Role
	claims["username"] = user.Username
	claims["uuid"] = user.Uuid
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString(MySigningKey)
	return tokenString
}

// Register a new user
func RegisterUser(user User) User {
	var checkUser User
	collection := client.Database("testing").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	docCount, err := collection.CountDocuments(ctx, User{})
	usernameErr := collection.FindOne(ctx, User{Username: user.Username}).Decode(&checkUser)
	emailErr := collection.FindOne(ctx, User{Email: user.Email}).Decode(&checkUser)
	if usernameErr == nil {
		return User{Username: "Taken"}
	} else if emailErr == nil {
		return User{Email: "Taken"}
	}
	curid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	user.Uuid = curid.String()
	password := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user.Password = string(hashedPassword)
	if docCount == 0 {
        user.Role = "admin"
	} else {
		user.Role = "member"
	}
	collection.InsertOne(ctx, user)
	return user
}

// Authenticate users signing in
func AuthenticateUser(user User) User {
	var checkUser User
	collection := client.Database("testing").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	usernameErr := collection.FindOne(ctx, User{Username: user.Username}).Decode(&checkUser)
	if usernameErr != nil {
		return User{Username: "NotFound"}
	}
	password := []byte(user.Password)
	checkPassword := []byte(checkUser.Password)
	err := bcrypt.CompareHashAndPassword(checkPassword, password)
	if err == nil {
		return checkUser
	} else {
		return User{Password: "Incorrect"}
	}
}

func RefreshToken() {

}

func DeleteUser() {

}

func FindUsers() {

}

func FindUser() {

}

/*
Todos Utilities
*/

// Find all Todos
func FindTodos() Todos {
	var todos Todos
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

// Find a specific to-do
func FindTodo(id string) Todo {
	var todo Todo
	collection := client.Database("testing").Collection("todos")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Todo{Uuid: id}).Decode(&todo)
	fmt.Println("---error:   --",err)
	fmt.Println("----todo: --", todo)
	if err != nil {
		return Todo{}
	}
	return todo
}

// Create new Todos
func CreateTodo(todo Todo) Todo {
	curid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	todo.Uuid = curid.String()
	collection := client.Database("testing").Collection("todos")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	collection.InsertOne(ctx, todo)
	return todo
}

// Delete todos
func DeleteTodo(id string) Todo {
	var todo Todo
	collection := client.Database("testing").Collection("todos")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOneAndDelete(ctx, Todo{Uuid: id}).Decode(&todo)
	if err != nil {
		return Todo{}
	}
	return todo
}
