/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/19/2019
*/


package main


import (
	"context"
	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"rest_api/pkg"
	"rest_api/pkg/configuration"
	"rest_api/pkg/mongodb"
	"rest_api/pkg/server"
	"time"
)


type App struct {
	server  *server.Server
	client  *mongo.Client
	config  configuration.Configuration
}



func (a *App) Initialize() {
	a.config = configuration.ConfigurationSettings()
	var err error
	a.client, err = mongodb.DatabaseConn(a.config.MongoURI)
	if err != nil {
		log.Fatalln("unable to connect to mongodb")
	}
	u := mongodb.NewUserService(a.client, a.config.Database, "users", a.config)
	g := mongodb.NewGroupService(a.client, a.config.Database, "groups")
	t := mongodb.NewTodoService(a.client, a.config.Database, "todos", a.config)
	// Create initial Admin User, Owner Group, and Bucket
	var group root.Group
	var adminUser root.User
	curid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	group.Name = a.config.DefaultAdminGroup
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	docCount, _ := a.client.Database(a.config.Database).Collection("groups").CountDocuments(ctx, bson.M{})
	if docCount == 0 {
		group.GroupType = "master_admins"
		group.Uuid = curid.String()
		adminGroup := g.GroupCreate(group)
		adminUser.Username = a.config.MasterAdminUsername
		adminUser.Email = a.config.MasterAdminEmail
		adminUser.Password = a.config.MasterAdminInitialPassword
		adminUser.GroupUuid = adminGroup.Uuid
		u.UserCreate(adminUser)
	}
	a.server = server.NewServer(u, g, t, a.config, a.client)

}
func (a *App) Run() {
	//defer a.client.Close()
	a.server.Start()
}
