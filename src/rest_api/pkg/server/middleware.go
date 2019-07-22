/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/19/2019
*/


package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"rest_api/pkg"
	"rest_api/pkg/configuration"
	"time"
)

type JWTError struct {
	Message string `json:"message"`
}

// Return JSON Error to Requested is Auth is bad
func respondWithError(w http.ResponseWriter, status int, error JWTError) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(error); err != nil {
		panic(err)
	}
}

// Middleware Function to Verify Requester is a Valid Admin
func AdminTokenVerifyMiddleWare(next http.HandlerFunc, config configuration.Configuration, client *mongo.Client) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var MySigningKey = []byte(config.Secret)
		var errorObject JWTError
		authToken := r.Header.Get("Auth-Token")
		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error")
			}
			return []byte(MySigningKey), nil
		})
		if err != nil {
			errorObject.Message = err.Error()
			respondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
		tokenClaims := token.Claims.(jwt.MapClaims)
		userUuid := tokenClaims["uuid"].(string)
		var checkUser root.User
		collection := client.Database(config.Database).Collection("users")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		usernameErr := collection.FindOne(ctx, bson.M{"uuid": userUuid}).Decode(&checkUser)
		if usernameErr != nil {
			errorObject.Message = "Invalid Token"
			respondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
		if token.Valid && checkUser.Role == "master_admin" {
			next.ServeHTTP(w, r)
		} else {
			errorObject.Message = "Invalid Token"
			respondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
	})
}

// Middleware Function to Verify Requested is Authenticated
func MemberTokenVerifyMiddleWare(next http.HandlerFunc, config configuration.Configuration, client *mongo.Client) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var MySigningKey = []byte(config.Secret)
		var errorObject JWTError
		authToken := r.Header.Get("Auth-Token")
		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error")
			}
			return []byte(MySigningKey), nil
		})
		if err != nil {
			errorObject.Message = err.Error()
			respondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
		tokenClaims := token.Claims.(jwt.MapClaims)
		userUuid := tokenClaims["uuid"].(string)
		var checkUser root.User
		collection := client.Database(config.Database).Collection("users")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		usernameErr := collection.FindOne(ctx, bson.M{"uuid": userUuid}).Decode(&checkUser)
		if usernameErr != nil {
			errorObject.Message = "Invalid Token"
			respondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			errorObject.Message = "Invalid Token"
			respondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
	})
}