/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/19/2019
*/

package server

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"io/ioutil"
	"net/http"
	"rest_api/pkg"
	"rest_api/pkg/configuration"
)

type groupRouter struct {
	groupService root.GroupService
}

// NewGroupRouter is a function that initializes a new groupRouter struct
func NewGroupRouter(g root.GroupService, router *mux.Router, config configuration.Configuration, client *mongo.Client) *mux.Router {
	groupRouter := groupRouter{g}
	router.HandleFunc("/groups", AdminTokenVerifyMiddleWare(groupRouter.GroupsShow, config, client)).Methods("GET")
	router.HandleFunc("/groups/{groupId}", AdminTokenVerifyMiddleWare(groupRouter.GroupShow, config, client)).Methods("GET")
	router.HandleFunc("/groups", AdminTokenVerifyMiddleWare(groupRouter.CreateGroup, config, client)).Methods("POST")
	router.HandleFunc("/groups/{groupId}", AdminTokenVerifyMiddleWare(groupRouter.DeleteGroup, config, client)).Methods("DELETE")
	router.HandleFunc("/groups/{groupId}", AdminTokenVerifyMiddleWare(groupRouter.ModifyGroup, config, client)).Methods("PATCH")
	return router
}

// Handler to show all groups
func (gr *groupRouter) ModifyGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupId := vars["groupId"]
	var group root.Group
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &group); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	group.Uuid = groupId
	g := gr.groupService.GroupUpdate(group)
	if g.Uuid == "Not Found" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(404)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Group Not Found"}); err != nil {
			panic(err)
		}
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusAccepted)
		if err := json.NewEncoder(w).Encode(g); err != nil {
			panic(err)
		}
	}
}

// Handler to show all groups
func (gr *groupRouter) GroupsShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	groups := gr.groupService.GroupsFind()
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		panic(err)
	}
}

// Handler to show a specific group
func (gr *groupRouter) GroupShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupId := vars["groupId"]
	group := gr.groupService.GroupFind(groupId)
	if group.Uuid != "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(group); err != nil {
			panic(err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Group Not Found"}); err != nil {
		panic(err)
	}
}

// Handler to create an group
func (gr *groupRouter) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group root.Group
	group.GroupType = "normal"
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &group); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	curid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	group.Uuid = curid.String()
	g := gr.groupService.GroupCreate(group)
	if g.Name == "Taken" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(403)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusForbidden, Text: "Group Name Taken"}); err != nil {
			panic(err)
		}
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(g); err != nil {
			panic(err)
		}
	}
}

// Handler to delete an group
func (gr *groupRouter) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupId := vars["groupId"]
	group := gr.groupService.GroupDelete(groupId)
	if group.Uuid != "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode("Group Deleted"); err != nil {
			panic(err)
		}
		return
	}
	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Group Not Found"}); err != nil {
		panic(err)
	}
}
