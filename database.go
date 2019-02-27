package main

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DatabaseConn() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(""))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	return client
}
