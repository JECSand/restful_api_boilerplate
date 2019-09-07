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

// UserService is used by the app to manage all user related controllers and functionality
type UserService struct {
	collection *mongo.Collection
	config     configuration.Configuration
	client     *mongo.Client
}

// NewUserService is an exported function used to initialize a new UserService struct
func NewUserService(client *mongo.Client, dbName string, collectionName string, config configuration.Configuration) *UserService {
	collection := client.Database(dbName).Collection(collectionName)
	return &UserService{collection, config, client}
}

// AuthenticateUser is used to authenticate users that are signing in
func (p *UserService) AuthenticateUser(user root.User) root.User {
	var checkUser = newUserModel(root.User{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
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
	}
	fmt.Println("password check err: ", err)
	return root.User{Password: "Incorrect"}
}

// BlacklistAuthToken is used during signout to add the now invalid auth-token/api key to the blacklist collection
func (p *UserService) BlacklistAuthToken(authToken string) {
	var blacklist root.Blacklist
	blacklist.AuthToken = authToken
	currentTime := time.Now().UTC()
	blacklist.LastModified = currentTime.String()
	blacklist.CreationDatetime = currentTime.String()
	blacklistModel := newBlacklistModel(blacklist)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	p.client.Database(p.config.Database).Collection("blacklists").InsertOne(ctx, blacklistModel)
}

// RefreshToken is used to refresh an existing & valid JWT token
func (p *UserService) RefreshToken(tokenData []string) root.User {
	if tokenData[0] == "" {
		return root.User{Uuid: ""}
	}
	userUuid := tokenData[0]
	user := p.UserFind(userUuid)
	return user
}

// UpdatePassword is used to update the currently logged in user's password
func (p *UserService) UpdatePassword(tokenData []string, CurrentPassword string, newPassword string) root.User {
	userUuid := tokenData[0]
	var user = newUserModel(root.User{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
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
	}
	return root.User{Password: "Incorrect"}
}

// UserCreate is used to create a new user
func (p *UserService) UserCreate(user root.User) root.User {
	var checkUser = newUserModel(root.User{})
	var checkGroup = newGroupModel(root.Group{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	docCount, err := p.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
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
	curid, err2 := uuid.NewV4()
	if err2 != nil {
		panic(err2)
	}
	user.Uuid = curid.String()
	password := []byte(user.Password)
	hashedPassword, err3 := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err3 != nil {
		panic(err3)
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

// UserDelete is used to delete an user
func (p *UserService) UserDelete(id string) root.User {
	var user = newUserModel(root.User{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := p.collection.FindOneAndDelete(ctx, bson.M{"uuid": id}).Decode(&user)
	if err != nil {
		return root.User{}
	}
	return user.toRootUser()
}

// UsersFind is used to find all user docs
func (p *UserService) UsersFind() []root.User {
	var users []root.User
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
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

// UserFind is used to find a specific user doc
func (p *UserService) UserFind(id string) root.User {
	var user = newUserModel(root.User{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := p.collection.FindOne(ctx, bson.M{"uuid": id}).Decode(&user)
	if err != nil {
		return root.User{}
	}
	return user.toRootUser()
}

// UserUpdate is used to update an existing user doc
func (p *UserService) UserUpdate(user root.User) root.User {
	var curUser = newUserModel(root.User{})
	var checkUser = newUserModel(root.User{})
	var checkGroup = newGroupModel(root.Group{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	docCount, err := p.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	userErr := p.collection.FindOne(ctx, bson.M{"uuid": user.Uuid}).Decode(&curUser)
	if userErr != nil {
		return root.User{Uuid: "Not Found"}
	}
	user = BaseModifyUser(user, curUser)
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
	}
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
				{"firstname", user.FirstName},
				{"lastname", user.LastName},
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
	}
	update := bson.D{{"$set",
		bson.D{
			{"firstname", user.FirstName},
			{"lastname", user.LastName},
			{"username", user.Username},
			{"email", user.Email},
			{"role", user.Role},
			{"groupuuid", user.GroupUuid},
			{"last_modified", currentTime.String()},
		},
	}}
	_, errUp := p.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("update err: ", errUp)
		panic(errUp)
	}
	return user
}

// UserDocInsert is used to insert an user doc directly into mongodb for testing purposes
func (p *UserService) UserDocInsert(user root.User) root.User {
	password := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user.Password = string(hashedPassword)
	var insertUser = newUserModel(user)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err2 := p.collection.InsertOne(ctx, insertUser)
	if err2 != nil {
		fmt.Println("todo doc insertion error: ", err2)
	}
	return insertUser.toRootUser()
}
