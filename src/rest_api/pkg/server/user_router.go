/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
9/05/2019
*/

package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"io/ioutil"
	"net/http"
	"rest_api/pkg"
	"rest_api/pkg/configuration"
	"time"
)

type userRouter struct {
	userService  root.UserService
	groupService root.GroupService
	config       configuration.Configuration
}

// NewUserRouter is a function that initializes a new userRouter struct
func NewUserRouter(u root.UserService, router *mux.Router, o root.GroupService, config configuration.Configuration, client *mongo.Client) *mux.Router {
	userRouter := userRouter{u, o, config}
	router.HandleFunc("/auth", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/auth", userRouter.Signin).Methods("POST")
	router.HandleFunc("/auth", MemberTokenVerifyMiddleWare(userRouter.RefreshSession, config, client)).Methods("GET")
	router.HandleFunc("/auth", MemberTokenVerifyMiddleWare(userRouter.Signout, config, client)).Methods("DELETE")
	router.HandleFunc("/auth/api-key", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/auth/api-key", MemberTokenVerifyMiddleWare(userRouter.GenerateAPIKey, config, client)).Methods("GET")
	router.HandleFunc("/auth/password", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/auth/password", MemberTokenVerifyMiddleWare(userRouter.UpdatePassword, config, client)).Methods("POST")
	router.HandleFunc("/users", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/users", AdminTokenVerifyMiddleWare(userRouter.UsersShow, config, client)).Methods("GET")
	router.HandleFunc("/users/{userId}", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/users/{userId}", AdminTokenVerifyMiddleWare(userRouter.UserShow, config, client)).Methods("GET")
	router.HandleFunc("/users", AdminTokenVerifyMiddleWare(userRouter.CreateUser, config, client)).Methods("POST")
	router.HandleFunc("/users/{userId}", AdminTokenVerifyMiddleWare(userRouter.DeleteUser, config, client)).Methods("DELETE")
	router.HandleFunc("/users/{userId}", AdminTokenVerifyMiddleWare(userRouter.ModifyUser, config, client)).Methods("PATCH")
	return router
}

// Handler function that manages the user signin process
func (ur *userRouter) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	decodedToken := DecodeJWT(r.Header.Get("Auth-Token"), ur.config)
	type passwordStruct struct {
		NewPassword     string `json:"new_password"`
		CurrentPassword string `json:"current_password"`
	}
	var pw passwordStruct
	err2 := json.Unmarshal(body, &pw)
	if err2 != nil {
		panic(err)
	}
	u := ur.userService.UpdatePassword(decodedToken, pw.CurrentPassword, pw.NewPassword)
	if u.Password == "Incorrect" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(403)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusForbidden, Text: "Incorrect Current Password Provided"}); err != nil {
			panic(err)
		}
	} else {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusAccepted)
		u.Password = ""
		if err := json.NewEncoder(w).Encode(u); err != nil {
			panic(err)
		}
	}
}

// Handler function that manages the user signin process
func (ur *userRouter) ModifyUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	var user root.User
	user.Uuid = userId
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &user); err != nil {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	u := ur.userService.UserUpdate(user)
	if u.Uuid == "Not Found" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(404)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "User Not Found"}); err != nil {
			panic(err)
		}
	}
	if u.Email == "Taken" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(403)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusForbidden, Text: "Email Taken"}); err != nil {
			panic(err)
		}
	} else if u.Username == "Taken" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(403)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusForbidden, Text: "Username Taken"}); err != nil {
			panic(err)
		}
	} else if u.GroupUuid == "No User Group Found" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(403)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Owner Not Found"}); err != nil {
			panic(err)
		}
	} else {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusAccepted)
		if err := json.NewEncoder(w).Encode(u); err != nil {
			panic(err)
		}
	}
}

// Handler function that manages the user signin process
func (ur *userRouter) Signin(w http.ResponseWriter, r *http.Request) {
	var user root.User
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &user); err != nil {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	u := ur.userService.AuthenticateUser(user)
	if u.Username == "NotFound" || u.Password == "Incorrect" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(401)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusUnauthorized, Text: "Incorrect"}); err != nil {
			panic(err)
		}
	} else {
		expDT := time.Now().Add(time.Hour * 1).Unix()
		sessionToken := CreateToken(u, ur.config, expDT)
		w = SetResponseHeaders(w, sessionToken, "")
		w.WriteHeader(http.StatusOK)
		u.Password = ""
		if err := json.NewEncoder(w).Encode(u); err != nil {
			panic(err)
		}
	}
}

// Handler function that refreshes a users JWT token
func (ur *userRouter) RefreshSession(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Auth-Token")
	tokenData := DecodeJWT(authToken, ur.config)
	user := ur.userService.RefreshToken(tokenData)
	expDT := time.Now().Add(time.Hour * 1).Unix()
	newToken := CreateToken(user, ur.config, expDT)
	w = SetResponseHeaders(w, newToken, "")
	w.WriteHeader(http.StatusOK)
}

// Handler function that generates a 6 month API Key for a given user
func (ur *userRouter) GenerateAPIKey(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Auth-Token")
	tokenData := DecodeJWT(authToken, ur.config)
	user := ur.userService.RefreshToken(tokenData)
	expDT := time.Now().Add(time.Hour * 4380).Unix()
	apiKey := CreateToken(user, ur.config, expDT)
	w = SetResponseHeaders(w, "", apiKey)
	w.WriteHeader(http.StatusOK)
}

// Handler function that ends a users session
func (ur *userRouter) Signout(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Auth-Token")
	ur.userService.BlacklistAuthToken(authToken)
	w = SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
}

// Handler function that creates a new user
func (ur *userRouter) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user root.User
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &user); err != nil {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	u := ur.userService.UserCreate(user)
	if u.Email == "Taken" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(403)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusForbidden, Text: "Email Taken"}); err != nil {
			panic(err)
		}
	} else if u.Username == "Taken" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(403)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusForbidden, Text: "Username Taken"}); err != nil {
			panic(err)
		}
	} else if u.GroupUuid == "No User Group Found" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(403)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusForbidden, Text: "Owner Not Found"}); err != nil {
			panic(err)
		}
	} else {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusCreated)
		u.Password = ""
		if err := json.NewEncoder(w).Encode(u); err != nil {
			panic(err)
		}
	}
}

// Handler that shows a specific user
func (ur *userRouter) UsersShow(w http.ResponseWriter, r *http.Request) {
	w = SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	users := ur.userService.UsersFind()
	if err := json.NewEncoder(w).Encode(users); err != nil {
		panic(err)
	}
}

// Handler that shows all users
func (ur *userRouter) UserShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	user := ur.userService.UserFind(userId)
	if user.Uuid != "" {
		user.Password = ""
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			panic(err)
		}
		return
	}
	w = SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "User Not Found"}); err != nil {
		panic(err)
	}
}

// Handler function that deletes a user
func (ur *userRouter) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	user := ur.userService.UserDelete(userId)
	if user.Uuid != "" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode("User Deleted"); err != nil {
			panic(err)
		}
		return
	}
	// If we didn't find it, 404
	w = SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "User Not Found"}); err != nil {
		panic(err)
	}
}
