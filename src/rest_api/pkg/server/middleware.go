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

// JWTError is a struct that is used to contain a json encoded error message for any JWT related errors
type JWTError struct {
	Message string `json:"message"`
}

// Return JSON Error to Requested is Auth is bad
func respondWithError(w http.ResponseWriter, status int, error JWTError) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Auth-Token")
	w.Header().Add("Access-Control-Expose-Headers", "Content-Type, Auth-Token")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(error); err != nil {
		panic(err)
	}
}

// Check submitted Auth-Token or API-Key with what's in the blacklist collection
func checkTokenBlacklist(authToken string, config configuration.Configuration, client *mongo.Client) bool {
	var checkToken root.Blacklist
	collection := client.Database(config.Database).Collection("blacklists")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	blacklistErr := collection.FindOne(ctx, bson.M{"auth_token": authToken}).Decode(&checkToken)
	if blacklistErr != nil {
		return false
	}
	return true
}

// AdminTokenVerifyMiddleWare is used to verify that the requester is a valid admin
func AdminTokenVerifyMiddleWare(next http.HandlerFunc, config configuration.Configuration, client *mongo.Client) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var MySigningKey = []byte(config.Secret)
		var errorObject JWTError
		authToken := r.Header.Get("Auth-Token")
		blacklistedToken := checkTokenBlacklist(authToken, config, client)
		if blacklistedToken {
			errorObject.Message = "Invalid Token"
			respondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
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
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		usernameErr := collection.FindOne(ctx, bson.M{"uuid": userUuid}).Decode(&checkUser)
		if usernameErr != nil {
			errorObject.Message = "Invalid Token"
			respondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
		if token.Valid && checkUser.Role == "master_admin" || checkUser.Role == "group_admin" {
			next.ServeHTTP(w, r)
		} else {
			errorObject.Message = "Invalid Token"
			respondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
	})
}

// MemberTokenVerifyMiddleWare is used to verify that a requester is authenticated
func MemberTokenVerifyMiddleWare(next http.HandlerFunc, config configuration.Configuration, client *mongo.Client) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var MySigningKey = []byte(config.Secret)
		var errorObject JWTError
		authToken := r.Header.Get("Auth-Token")
		blacklistedToken := checkTokenBlacklist(authToken, config, client)
		if blacklistedToken {
			errorObject.Message = "Invalid Token"
			respondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
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
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
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
