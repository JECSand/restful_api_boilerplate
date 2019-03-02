package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/dgrijalva/jwt-go"
)

type JWTError struct {
	Message string `json:"message"`
}

func respondWithError(w http.ResponseWriter, status int, error JWTError) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(error); err != nil {
		panic(err)
	}
}


func AdminTokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		var checkUser User
		collection := client.Database("testing").Collection("users")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		usernameErr := collection.FindOne(ctx, User{Uuid: userUuid}).Decode(&checkUser)
		if usernameErr != nil {
			errorObject.Message = "Invalid Token"
			respondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
		if token.Valid && checkUser.Role == "admin" {
			next.ServeHTTP(w, r)
		} else {
			errorObject.Message = "Invalid Token"
			respondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
	})
}


func MemberTokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		var checkUser User
		collection := client.Database("testing").Collection("users")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		usernameErr := collection.FindOne(ctx, User{Uuid: userUuid}).Decode(&checkUser)
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