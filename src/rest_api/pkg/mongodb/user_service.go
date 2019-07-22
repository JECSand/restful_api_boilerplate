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
	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"rest_api/pkg"
	"rest_api/pkg/configuration"
	"time"
)


type UserService struct {
	collection *mongo.Collection
	config     configuration.Configuration
	client     *mongo.Client
}



func NewUserService(client *mongo.Client, dbName string, collectionName string, config configuration.Configuration) *UserService {
	collection := client.Database(dbName).Collection(collectionName)
	return &UserService {collection, config, client}
}


// Authenticate users signing in
func (p *UserService) AuthenticateUser(user root.User) root.User {
	var checkUser = newUserModel(root.User{})
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	usernameErr := p.collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&checkUser)
	if usernameErr != nil {
		fmt.Println("Look up error: ", usernameErr)
		return root.User{Username: "NotFound"}
	}
	rootUser := checkUser.toRootUser()
	password := []byte(user.Password)
	checkPassword := []byte(checkUser.Password)
	err := bcrypt.CompareHashAndPassword(checkPassword, password)
	if err == nil {
		return rootUser
	} else {
		fmt.Println("password check err: ", err)
		return root.User{Password: "Incorrect"}
	}
}


// Refresh Existing Session JWT Token
func (p *UserService) RefreshToken(tokenData []string) root.User {
	if tokenData[0] == "" {
		return root.User{Uuid: ""}
	}
	userUuid := tokenData[0]
	user := p.UserFind(userUuid)
	return user
}


// Update your user's password
func (p *UserService) UpdatePassword(tokenData []string, CurrentPassword string, newPassword string) root.User {
	userUuid := tokenData[0]
	var user = newUserModel(root.User{})
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := p.collection.FindOne(ctx, bson.M{"uuid": userUuid}).Decode(&user)
	if err != nil {
		return root.User{Uuid: "Not Found"}
	}
	// 2. Check current password
	curUser := user.toRootUser()
	password := []byte(CurrentPassword)
	checkPassword := []byte(curUser.Password)
	err2 := bcrypt.CompareHashAndPassword(checkPassword, password)
	if err2 == nil {
		// 3. Update doc with new password
		currentTime := time.Now().UTC()
		hashedPassword, err3 := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err3 != nil {
			panic(err3)
		}
		filter := bson.D{{"uuid", curUser.Uuid}}
		update := bson.D{{"$set",
			bson.D{
				{"password", string(hashedPassword)},
				{"last_modified", currentTime.String()},
			},
		}}
		_, err4 := p.collection.UpdateOne(ctx, filter, update)
		if err4 != nil {
			fmt.Println("update err: ", err4)
			panic(err4)
		}
		return curUser
	} else {
		return root.User{Password: "Incorrect"}
	}
	return root.User{}
}


// Create a new user
func (p *UserService) UserCreate(user root.User) root.User {
	var checkUser = newUserModel(root.User{})
	var checkGroup = newGroupModel(root.Group{})
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	docCount, err := p.collection.CountDocuments(ctx, bson.M{})
	usernameErr := p.collection.FindOne(ctx, bson.M{"username": user.Username, "groupuuid": user.GroupUuid}).Decode(&checkUser)
	emailErr := p.collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&checkUser)
	groupErr := p.client.Database(p.config.Database).Collection("groups").FindOne(ctx, bson.M{"uuid": user.GroupUuid}).Decode(&checkGroup)
	if usernameErr == nil {
		return root.User{Username: "Taken"}
	} else if emailErr == nil {
		return root.User{Email: "Taken"}
	} else if groupErr != nil {
		return root.User{GroupUuid: "No User Group Found"}
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
		user.Role = "master_admin"
	} else {
		if user.Role != "master_admin" {
			user.Role = "member"
		} else {
			var masterGroup = newGroupModel(root.Group{})
			groupErr := p.client.Database(p.config.Database).Collection("groups").FindOne(ctx, bson.M{"name": p.config.DefaultAdminGroup}).Decode(&masterGroup)
			if groupErr != nil {
				return root.User{GroupUuid: "No User Group Found"}
			}
			user.GroupUuid = masterGroup.Uuid
		}
	}
	currentTime := time.Now().UTC()
	user.LastModified = currentTime.String()
	user.CreationDatetime = currentTime.String()
	userModel := newUserModel(user)
	p.collection.InsertOne(ctx, userModel)
	return userModel.toRootUser()
}


// Delete a User
func (p *UserService) UserDelete(id string) root.User {
	var user = newUserModel(root.User{})
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := p.collection.FindOneAndDelete(ctx, bson.M{"uuid": id}).Decode(&user)
	if err != nil {
		return root.User{}
	}
	return user.toRootUser()
}


// Find Users
func (p *UserService) UsersFind() []root.User {
	var users []root.User
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := p.collection.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user = newUserModel(root.User{})
		cursor.Decode(&user)
		user.Password = ""
		users = append(users, user.toRootUser())
	}
	return users
}


// Find a User
func (p *UserService) UserFind(id string) root.User {
	var user = newUserModel(root.User{})
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := p.collection.FindOne(ctx, bson.M{"uuid": id}).Decode(&user)
	if err != nil {
		return root.User{}
	}
	return user.toRootUser()
}


// Update an existing user
func (p *UserService) UserUpdate(user root.User) root.User {
	var curUser = newUserModel(root.User{})
	var checkUser = newUserModel(root.User{})
	var checkGroup = newGroupModel(root.Group{})
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	docCount, err := p.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	userErr := p.collection.FindOne(ctx, bson.M{"uuid": user.Uuid}).Decode(&curUser)
	if userErr != nil {
		return root.User{Uuid: "Not Found"}
	}
	if len(user.Username) == 0 { user.Username = curUser.Username }
	if len(user.Email) == 0 { user.Email = curUser.Email }
	if len(user.GroupUuid) == 0 { user.GroupUuid = curUser.GroupUuid }
	if len(user.Role) == 0 { user.Role = curUser.Role }
	usernameErr := p.collection.FindOne(ctx, bson.M{"username": user.Username, "groupuuid": user.GroupUuid}).Decode(&checkUser)
	emailErr := p.collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&checkUser)
	groupErr := p.client.Database(p.config.Database).Collection("groups").FindOne(ctx, bson.M{"uuid": user.GroupUuid}).Decode(&checkGroup)
	if usernameErr == nil {
		return root.User{Username: "Taken"}
	} else if emailErr == nil {
		return root.User{Email: "Taken"}
	} else if groupErr != nil {
		return root.User{GroupUuid: "No User Group Found"}
	}
	if docCount == 0 {
		return root.User{Uuid: "Not Found"}
	} else {
		filter := bson.D{{"uuid", user.Uuid}}
		currentTime := time.Now().UTC()
		if len(user.Password) != 0 {
			password := []byte(user.Password)
			hashedPassword, hashErr := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
			if hashErr != nil {
				panic(hashErr)
			}
			update := bson.D{{"$set",
				bson.D{
					{"password", string(hashedPassword)},
					{"username", user.Username},
					{"email", user.Email},
					{"role", user.Role},
					{"groupuuid", user.GroupUuid},
					{"last_modified", currentTime.String()},
				},
			}}
			_, err := p.collection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println("update err: ", err)
				panic(err)
			}
			user.Password = ""
			return user
		} else {
			update := bson.D{{"$set",
				bson.D{
					{"username", user.Username},
					{"email", user.Email},
					{"role", user.Role},
					{"groupuuid", user.GroupUuid},
					{"last_modified", currentTime.String()},
				},
			}}
			_, err := p.collection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println("update err: ", err)
				panic(err)
			}
			return user
		}
	}
	return root.User{}
}