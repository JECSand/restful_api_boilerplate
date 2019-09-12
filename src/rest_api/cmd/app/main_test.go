/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/25/2019
*/

package main

import (
	"bytes"
	"net/http"
	"testing"
)

var ta App

// Setup Tests
func setup() {
	ta = App{}
	ta.Initialize("test")
}

/*
AUTH & USER TESTS
*/

// User Signin Test
func TestSignin(t *testing.T) {
	setup()
	testResponse := signin(ta, "MasterAdmin", "123xyzabc")
	clean(ta)
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}

// Create User Test
func TestCreateUser(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	authResponse := signin(ta, "MasterAdmin", "123xyzabc")
	authToken := authResponse.Header().Get("Auth-Token")
	// Create User Test Request
	payload := []byte(`{"username":"test_user","password":"789test123","firstname":"test","lastname":"user","email":"test@example.com","groupuuid":"96840292-0dd5-44a8-a143-92d949a6be36","role":"member"}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean test database and check test response
	clean(ta)
	checkResponseCode(t, http.StatusCreated, testResponse.Code)
}

// User Signout Test
func TestSignout(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	authResponse := signin(ta, "MasterAdmin", "123xyzabc")
	authToken := authResponse.Header().Get("Auth-Token")
	// Signout Test
	reqSignout, _ := http.NewRequest("DELETE", "/auth", nil)
	reqSignout.Header.Add("Content-Type", "application/json")
	reqSignout.Header.Add("Auth-Token", authToken)
	signoutTestResponse := executeRequest(ta, reqSignout)
	checkResponseCode(t, http.StatusOK, signoutTestResponse.Code)
	// Test to ensure the Auth-Token or API-Key is now invalid
	payload := []byte(`{"username":"test_user","password":"789test123","firstname":"tester","lastname":"userington","email":"test@example.com","groupuuid":"96840292-0dd5-44a8-a143-92d949a6be36","role":"member"}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean test database and check test response
	clean(ta)
	checkResponseCode(t, http.StatusUnauthorized, testResponse.Code)
}

// Update Password Test
func TestUpdatePassword(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	createTestUser(ta, 1)
	authResponse := signin(ta, "test_user", "789test123")
	authToken := authResponse.Header().Get("Auth-Token")
	// Update User Password Test Request with incorrect current password
	payloadErr := []byte(`{"current_password":"789test122","new_password":"789test124"}`)
	reqErr, _ := http.NewRequest("POST", "/auth/password", bytes.NewBuffer(payloadErr))
	reqErr.Header.Add("Content-Type", "application/json")
	reqErr.Header.Add("Auth-Token", authToken)
	testResponseErr := executeRequest(ta, reqErr)
	checkResponseCode(t, http.StatusForbidden, testResponseErr.Code)
	// Test to ensure password did not update due to incorrect current password
	testResponseErrAuth := signin(ta, "test_user", "789test124")
	checkResponseCode(t, http.StatusUnauthorized, testResponseErrAuth.Code)
	// Update User Password Test Request correctly
	payloadOK := []byte(`{"current_password":"789test123","new_password":"789test124"}`)
	reqOK, _ := http.NewRequest("POST", "/auth/password", bytes.NewBuffer(payloadOK))
	reqOK.Header.Add("Content-Type", "application/json")
	reqOK.Header.Add("Auth-Token", authToken)
	testResponseOK := executeRequest(ta, reqOK)
	checkResponseCode(t, http.StatusAccepted, testResponseOK.Code)
	// Test to ensure user can now log in with new password
	testResponseOKAuth := signin(ta, "test_user", "789test124")
	// Clean database and do final status check
	clean(ta)
	checkResponseCode(t, http.StatusOK, testResponseOKAuth.Code)
}

// Modify User Test
func TestModifyUser(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	createTestGroup(ta, 2)
	createTestUser(ta, 1)
	authResponse := signin(ta, "MasterAdmin", "123xyzabc")
	authToken := authResponse.Header().Get("Auth-Token")
	// Modify User Document Test
	payloadUpdate := []byte(`{"username":"newUserName","password":"newUserpass","email":"new_test@email.com","groupuuid":"66846392-0ee5-13a8-a143-92d919a6be33","role":"member"}`)
	reqUpdate, _ := http.NewRequest("PATCH", "/users/dc20fd0a-eaad-43ca-9452-0a0b334cff8f", bytes.NewBuffer(payloadUpdate))
	reqUpdate.Header.Add("Content-Type", "application/json")
	reqUpdate.Header.Add("Auth-Token", authToken)
	updateTestResponse := executeRequest(ta, reqUpdate)
	checkResponseCode(t, http.StatusAccepted, updateTestResponse.Code)
	// Attempt to test Modified user doc by loggin in with new password/username
	testResponseOKAuth := signin(ta, "newUserName", "newUserpass")
	// Clean database and do final status check
	clean(ta)
	checkResponseCode(t, http.StatusOK, testResponseOKAuth.Code)
}

// List Users Test
func TestListUsers(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	createTestUser(ta, 1)
	createTestUser(ta, 2)
	authResponse := signin(ta, "MasterAdmin", "123xyzabc")
	authToken := authResponse.Header().Get("Auth-Token")
	// List all users test
	req, _ := http.NewRequest("GET", "/users", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	clean(ta)
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}

// List User Test
func TestListUser(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	createTestUser(ta, 1)
	authResponse := signin(ta, "MasterAdmin", "123xyzabc")
	authToken := authResponse.Header().Get("Auth-Token")
	// List a specific user Test
	req, _ := http.NewRequest("GET", "/users/dc20fd0a-eaad-43ca-9452-0a0b334cff8f", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	clean(ta)
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}

// Delete User Test
func TestDeleteUser(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	createTestUser(ta, 1)
	authResponse := signin(ta, "MasterAdmin", "123xyzabc")
	authToken := authResponse.Header().Get("Auth-Token")
	// Delete a user Test
	req, _ := http.NewRequest("DELETE", "/users/dc20fd0a-eaad-43ca-9452-0a0b334cff8f", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	clean(ta)
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}

// Refresh Auth Token Test
func TestTokenRefresh(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	authResponse := signin(ta, "MasterAdmin", "123xyzabc")
	authToken := authResponse.Header().Get("Auth-Token")
	// Refresh Auth-Token Test
	reqRefresh, _ := http.NewRequest("GET", "/auth", nil)
	reqRefresh.Header.Add("Content-Type", "application/json")
	reqRefresh.Header.Add("Auth-Token", authToken)
	refreshTestResponse := executeRequest(ta, reqRefresh)
	checkResponseCode(t, http.StatusOK, refreshTestResponse.Code)
	authToken = refreshTestResponse.Header().Get("Auth-Token")
	// Test new token to ensure it works by creating a new group
	payload := []byte(`{"username":"test_user","password":"789test123","email":"test@example.com","groupuuid":"96840292-0dd5-44a8-a143-92d949a6be36","role":"member"}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	clean(ta)
	checkResponseCode(t, http.StatusCreated, testResponse.Code)
}

// Generate API Key Test
func TestGenerateAPIKey(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	authResponse := signin(ta, "MasterAdmin", "123xyzabc")
	authToken := authResponse.Header().Get("Auth-Token")
	// Refresh Auth-Token Test
	reqAPIKey, _ := http.NewRequest("GET", "/auth/api-key", nil)
	reqAPIKey.Header.Add("Content-Type", "application/json")
	reqAPIKey.Header.Add("Auth-Token", authToken)
	apiKeyTestResponse := executeRequest(ta, reqAPIKey)
	checkResponseCode(t, http.StatusOK, apiKeyTestResponse.Code)
	apiKey := apiKeyTestResponse.Header().Get("API-Key")
	// Test new token to ensure it works by creating a new group
	payload := []byte(`{"username":"test_user","password":"789test123","email":"test@example.com","groupuuid":"96840292-0dd5-44a8-a143-92d949a6be36","role":"member"}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", apiKey)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	clean(ta)
	checkResponseCode(t, http.StatusCreated, testResponse.Code)
}

/*
GROUP TESTS
*/

// Create Group Test
func TestCreateGroup(t *testing.T) {
	// Test Setup
	setup()
	authResponse := signin(ta, "MasterAdmin", "123xyzabc")
	authToken := authResponse.Header().Get("Auth-Token")
	// Create new group test
	payload := []byte(`{"name":"testingGroup"}`)
	req, _ := http.NewRequest("POST", "/groups", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	clean(ta)
	checkResponseCode(t, http.StatusCreated, testResponse.Code)
}

// Modify Group Test
func TestModifyGroup(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	authResponse := signin(ta, "MasterAdmin", "123xyzabc")
	authToken := authResponse.Header().Get("Auth-Token")
	// Modify group test
	payload := []byte(`{"name":"newTestingGroup"}`)
	req, _ := http.NewRequest("PATCH", "/groups/96840292-0dd5-44a8-a143-92d949a6be36", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	clean(ta)
	checkResponseCode(t, http.StatusAccepted, testResponse.Code)
}

// List Groups Test
func TestListGroups(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	createTestGroup(ta, 2)
	authResponse := signin(ta, "MasterAdmin", "123xyzabc")
	authToken := authResponse.Header().Get("Auth-Token")
	// List all groups test
	req, _ := http.NewRequest("GET", "/groups", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	clean(ta)
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}

// List Group Test
func TestListGroup(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	authResponse := signin(ta, "MasterAdmin", "123xyzabc")
	authToken := authResponse.Header().Get("Auth-Token")
	// List all groups test
	req, _ := http.NewRequest("GET", "/groups/96840292-0dd5-44a8-a143-92d949a6be36", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	clean(ta)
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}

// Delete Group Test
func TestDeleteGroup(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	authResponse := signin(ta, "MasterAdmin", "123xyzabc")
	authToken := authResponse.Header().Get("Auth-Token")
	// List all groups test
	req, _ := http.NewRequest("DELETE", "/groups/96840292-0dd5-44a8-a143-92d949a6be36", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	clean(ta)
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}

/*
TODOS TESTS
*/

// Create Todos Test
func TestCreateTodo(t *testing.T) {
	// Test Setup
	setup()
	authResponse := signin(ta, "MasterAdmin", "123xyzabc")
	authToken := authResponse.Header().Get("Auth-Token")
	// Create new todos test
	payload := []byte(`{"name":"test_normal","due":"2019-01-02 15:04:11 -0700 UTC","description":"test@example.com"}`)
	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	clean(ta)
	checkResponseCode(t, http.StatusCreated, testResponse.Code)
}

// Modify Todos Test
func TestModifyTodo(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	createTestUser(ta, 1)
	createTestTodo(ta, 1)
	authResponse := signin(ta, "test_user", "789test123")
	authToken := authResponse.Header().Get("Auth-Token")
	// Modify todos test
	payload := []byte(`{"name":"new_todo_name","due":"2019-08-06 12:04:01 -0000 UTC","description":"Updated Task to complete","completed":"true"}`)
	req, _ := http.NewRequest("PATCH", "/todos/2e4c630e-89c8-4830-9841-f1b090fedc41", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	clean(ta)
	checkResponseCode(t, http.StatusAccepted, testResponse.Code)
}

// List Todos Test
func TestListTodos(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	createTestUser(ta, 1)
	createTestTodo(ta, 1)
	createTestTodo(ta, 2)
	authResponse := signin(ta, "test_user", "789test123")
	authToken := authResponse.Header().Get("Auth-Token")
	// List all todos test
	req, _ := http.NewRequest("GET", "/todos", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	clean(ta)
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}

// List Todos Test
func TestListTodo(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	createTestUser(ta, 1)
	createTestTodo(ta, 1)
	authResponse := signin(ta, "test_user", "789test123")
	authToken := authResponse.Header().Get("Auth-Token")
	// List a specific todos doc
	req, _ := http.NewRequest("GET", "/todos/2e4c630e-89c8-4830-9841-f1b090fedc41", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	clean(ta)
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}

// Delete Todos Test
func TestDeleteTodo(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	createTestUser(ta, 1)
	createTestTodo(ta, 1)
	authResponse := signin(ta, "test_user", "789test123")
	authToken := authResponse.Header().Get("Auth-Token")
	// List a specific todos doc
	req, _ := http.NewRequest("DELETE", "/todos/2e4c630e-89c8-4830-9841-f1b090fedc41", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	clean(ta)
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}