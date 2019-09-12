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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// DatabaseConn is a function that takes a mongoUri string and outputs a connected mongo client for the app to use
func DatabaseConn(mongoUri string) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		fmt.Println("Mongodb Connection Error!")
	}
	return client, err
}
