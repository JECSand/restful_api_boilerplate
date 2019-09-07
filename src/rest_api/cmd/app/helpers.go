/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/19/2019
*/

package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"rest_api/pkg"
	"testing"
	"time"
)

// Execute test an http request
func executeRequest(ta App, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	ta.server.Router.ServeHTTP(rr, req)
	return rr
}

// Check response code returned from a test http request
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

// Signin
func signin(ta App, username string, password string) *httptest.ResponseRecorder {
	payload := []byte(`{"username":"` + username + `","password":"` + password + `"}`)
	req, _ := http.NewRequest("POST", "/auth", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	response := executeRequest(ta, req)
	return response
}

// Create a group doc for test setup
func createTestGroup(ta App, groupType int) root.Group {
	group := root.Group{}
	if groupType == 1 {
		group.GroupType = "normal"
		group.Uuid = "96840292-0dd5-44a8-a143-92d949a6be36"
		group.Name = "newGroupName"
		group.LastModified = "2019-07-22 05:36:20.751969156 +0000 UTC"
		group.CreationDatetime = "2019-07-22 05:11:29.922897535 +0000 UTC"
	} else {
		group.GroupType = "normal"
		group.Uuid = "66846392-0ee5-13a8-a143-92d919a6be33"
		group.Name = "Group2"
		group.LastModified = "2019-07-22 05:36:20.751969156 +0000 UTC"
		group.CreationDatetime = "2019-07-22 05:11:29.922897535 +0000 UTC"
	}
	ta.server.GroupService.GroupDocInsert(group)
	return group
}

// Create a user doc for test setup
func createTestUser(ta App, userType int) root.User {
	user := root.User{}
	if userType == 1 {
		user.Uuid = "dc20fd0a-eaad-43ca-9452-0a0b334cff8f"
		user.Username = "test_user"
		user.Password = "789test123"
		user.FirstName = "test"
		user.LastName = "user"
		user.Email = "test@example.com"
		user.Role = "member"
		user.GroupUuid = "96840292-0dd5-44a8-a143-92d949a6be36"
		user.LastModified = "2019-07-22 05:36:20.751969156 +0000 UTC"
		user.CreationDatetime = "2019-07-22 05:11:29.922897535 +0000 UTC"
	} else {
		user.Uuid = "ex21fd0a-ebbb-73cv-2082-0x1b944cie9f"
		user.Username = "test_user2"
		user.Password = "789test124"
		user.FirstName = "tester"
		user.LastName = "userington"
		user.Email = "test2@example.com"
		user.Role = "member"
		user.GroupUuid = "96840292-0dd5-44a8-a143-92d949a6be36"
		user.LastModified = "2019-07-22 05:36:20.751969156 +0000 UTC"
		user.CreationDatetime = "2019-07-22 05:11:29.922897535 +0000 UTC"
	}
	ta.server.UserService.UserDocInsert(user)
	return user
}

// Create a todos doc for test setup
func createTestTodo(ta App, todoType int) root.Todo {
	todo := root.Todo{}
	if todoType == 1 {
		todo.Uuid = "2e4c630e-89c8-4830-9841-f1b090fedc41"
		todo.Name = "todo_name"
		todo.Completed = "false"
		todo.Due = "2019-08-06 12:04:01 +0000 UTC"
		todo.Description = "Updated Task to complete"
		todo.UserUuid = "dc20fd0a-eaad-43ca-9452-0a0b334cff8f"
		todo.GroupUuid = "96840292-0dd5-44a8-a143-92d949a6be36"
		todo.LastModified = "2019-07-22 05:29:53.414952125 +0000 UTC"
		todo.CreationDatetime = "2019-07-22 05:28:21.280285025 +0000 UTC"
	} else {
		todo.Uuid = "0e2c165d-83c2-4bc9-98c6-1402e589d1d2"
		todo.Name = "todo_name2"
		todo.Completed = "false"
		todo.Due = "2019-08-08 12:04:01 +0000 UTC"
		todo.Description = "Updated Task to complete2"
		todo.UserUuid = "dc20fd0a-eaad-43ca-9452-0a0b334cff8f"
		todo.GroupUuid = "96840292-0dd5-44a8-a143-92d949a6be36"
		todo.LastModified = "2019-07-23 05:29:53.414952125 +0000 UTC"
		todo.CreationDatetime = "2019-07-23 05:28:21.280285025 +0000 UTC"
	}
	ta.server.TodoService.TodoDocInsert(todo)
	return todo
}

// Cleanup after a test
func clean(ta App) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ta.client.Database("test").Drop(ctx)
}

