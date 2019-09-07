# Go RESTful API Boilerplate

A RESTful API Boilerplate written in Go.


[![Build Status](https://travis-ci.org/JECSand/restful_api_boilerplate.svg?branch=master)](https://travis-ci.org/JECSand/restful_api_boilerplate)
[![Go Report Card](https://goreportcard.com/badge/github.com/JECSand/restful_api_boilerplate)](https://goreportcard.com/report/github.com/JECSand/restful_api_boilerplate)


* Author(s): John Connor Sanders
* Current Version: 0.5.0
* Release Date: 7/19/2019
* MIT License
___
## Getting Started

Follow the instructions below to get the Go RESTful API up and running on your Linux Environment

### Prerequisites


* An Ubuntu 18+ or CentOS 7+ Operating System
* MongoDB 4+

### Setup

1. Clone the git repo:
```bash
$ git clone https://github.com/JECSand/restful_api_boilerplate.git
$ cd restful_api_boilerplate
```

2. Install requirements using the install.sh script:

```bash
$ . ./install.sh
```

3. Create conf.json file from the example file and configure the following settings:

* MongoDB URI Connection String
* Secret Encryption String
* Mongo Database Name
* Master Admin Username
* Master Admin Email
* Master Admin Initial Password
* Layout of the DateTime Strings (Can be left as default)
* Whether to run App with HTTPS
* If HTTPS is on, the cert.pem file
* If HTTPS is on, the path to the key.pem file
* Whether you want new users to be able to sign themselves up for accounts

```bash
$ cp conf.json.example conf.json
$ vi conf.json
```
___
## Running the API

### Production

To start the API's production build, run the following:

```bash
$ sh start.sh
```

To stop the API:

```bash
$ sh stop.sh
```

### Development

To start the API in development:

```bash
$ go build ./src/rest_api/cmd/app
$ ./app
```

To stop the development API, enter 'ctrl + c'

### Testing

To start the API's testing module, run the following:

```bash
$ go test ./src/rest_api/cmd/app
```

___
## API Route Guide
### I) Authentication Routes

___
#### 1. Signin
* POST - /auth

##### Request

***
* Headers

```
{
  Content-Type: application/json
}
```

* Body
```
{
  "username": "userName",
  "password": "userpass"
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Auth-Token: "",
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH
}
```

* Body
```
{
  "id": "000000000000000000000000",
  "uuid": "00000000-0000-0000-0000-000000000000",
  "username": "userName",
  "firstname": "john",
  "lastname": "smith",
  "email": "user@example.com",
  "role": "member",
  "groupuuid": "00000000-0000-0000-0000-000000000000",
  "last_modified": "2019-06-07 20:17:14.630917778 +0000 UTC",
  "creation_datetime": "2019-06-07 20:17:14.630917778 +0000 UTC"
}
```

#### 2. Signup
* POST - /auth/register
* This route will return a 404 if the "Registration" setting is set to "off" in the conf.json file.

##### Request

***
* Headers

```
{
  Content-Type: application/json
}
```

* Body
```
{
  "firstname": "john",
  "lastname": "smith",
  "email": "user@example.com"",
  "username": "userName",
  "password": "userpass"
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Auth-Token: "",
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

* Body
```
{
  "id": "000000000000000000000000",
  "uuid": "00000000-0000-0000-0000-000000000000",
  "username": "userName",
  "firstname": "john",
  "lastname": "smith",
  "email": "user@example.com",
  "role": "member",
  "groupuuid": "00000000-0000-0000-0000-000000000000",
  "last_modified": "2019-06-07 20:17:14.630917778 +0000 UTC",
  "creation_datetime": "2019-06-07 20:17:14.630917778 +0000 UTC"
}
```

#### 3. Refresh Token
* GET - /auth

##### Request

***
* Headers

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Auth-Token: "",
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

#### 4. Signout
* DELETE - /auth

##### Request

***
* Headers

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH
}
```

#### 5. API Key - Expires after 6 Months
* GET - /auth/api-key

##### Request

***
* Headers

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Auth-Token: "",
  API-Key: "",
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

#### 6. Update Password
* POST - /auth/password

##### Request

***
* Headers

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

***
* Body

```
{
  "current_password": "current_password",
  "new_password": "new_password"
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

### II) Todo Routes

___
#### 1. List Todo(s)
* GET - /todos/{todoId}
* todoId parameter is optional, if used the request will only return an object for that item.

##### Request

***
* Headers
```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

##### Response

***
* Headers
```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

* Body
```
[
  {
    "id":  "000000000000000000000000",
    "uuid": "00000000-0000-0000-0000-000000000000",
    "name": "todo_name",
    "due": "2019-08-01 12:04:01 -0000 UTC",
    "description": "Task to complete",
    "useruuid": "00000000-0000-0000-0000-000000000000",
    "groupuuid": "00000000-0000-0000-0000-000000000000",
    "last_modified": "2019-06-07 20:28:09.400248747 +0000 UTC",
    "creation_datetime": "2019-06-07 20:28:09.400248747 +0000 UTC"
  }
]
```

#### 2. Create Todo
* POST - /todos

##### Request

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Auth-Token: ""
}
```

* Body
```
{
    "name": "todo_name",
    "due": "2019-08-01 12:04:01 -0000 UTC",
    "description": "Task to complete"
}
```


##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

* Body
```
{
  "id":  "000000000000000000000000",
  "uuid": "00000000-0000-0000-0000-000000000000",
  "name": "todo_name",
  "completed": "false",
  "due": "2019-08-01 12:04:01 -0000 UTC",
  "description": "Task to complete",
  "useruuid": "00000000-0000-0000-0000-000000000000",
  "groupuuid": "00000000-0000-0000-0000-000000000000",
  "last_modified": "2019-06-07 20:28:09.400248747 +0000 UTC",
  "creation_datetime": "2019-06-07 20:28:09.400248747 +0000 UTC"
}
```

#### 3. Modify Todo
* PATCH - /todos/{todoId}

##### Request

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Auth-Token: ""
}
```

* Body
```
{
    "name": "new_todo_name",
    "due": "2019-08-06 12:04:01 -0000 UTC",
    "description": "Updated Task to complete",
    "completed": "true",
    "useruuid": "00000000-0000-0000-0000-000000000000"
}
```


##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

* Body
```
{
  "id":  "000000000000000000000000",
  "uuid": "00000000-0000-0000-0000-000000000000",
  "name": "new_todo_name",
  "completed": "true",
  "due": "2019-08-06 12:04:01 -0000 UTC",
  "description": "Task to complete",
  "useruuid": "00000000-0000-0000-0000-000000000000",
  "groupuuid": "00000000-0000-0000-0000-000000000000",
  "last_modified": "2019-06-07 20:28:09.400248747 +0000 UTC",
  "creation_datetime": "2019-06-07 20:28:09.400248747 +0000 UTC"
}
```

#### 4. Delete Todo
* DELETE - /todos/{todoId}

##### Request

***
* Headers

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

### III) Users Routes (Admins Only)

___
#### 1. List User(s)
* GET - /users/{userId}
* userId parameter is optional, if used the request will only return an object for that item.
   
##### Request

***
* Headers
   
```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```
   
##### Response

***
* Headers
   
```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```
   
* Body
```
[
  {
    "id": "000000000000000000000000",
    "uuid": "00000000-0000-0000-0000-000000000000",
    "username": "userName",
    "firstname": "jane",
    "lastname": "smith",
    "email": "user@example.com",
    "role": "member",
    "groupuuid": "00000000-0000-0000-0000-000000000000",
    "last_modified": "2019-06-07 20:17:14.630917778 +0000 UTC",
    "creation_datetime": "2019-06-07 20:17:14.630917778 +0000 UTC"
  }
]
```
   
   
#### 2. Create User
* POST - /users
  
##### Request

***
* Headers
   
```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```
   
* Body
```
{
  "username": "userName",
  "password": "userpass",
  "firstname": "jane",
  "lastname": "smith",
  "email": "test@email.com",
  "groupuuid": "00000000-0000-0000-0000-000000000000",
  "role": "member"
}
```
   
##### Response

***
* Headers
   
```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```
   
* Body
```
{
  "id": "000000000000000000000000",
  "uuid": "00000000-0000-0000-0000-000000000000",
  "username": "userName",
  "firstname": "jane",
  "lastname": "smith",
  "email": "test@email.com",
  "role": "member",
  "groupuuid": "00000000-0000-0000-0000-000000000000",
  "last_modified": "2019-06-07 20:22:44.322303274 +0000 UTC",
  "creation_datetime": "2019-06-07 20:22:44.322303274 +0000 UTC"
}
```

#### 3. Modify User
* PATCH - /users/{userId}
  
##### Request

***
* Headers
   
```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```
   
* Body
```
{
  "username": "newUserName",
  "password": "newUserpass",
  "firstname": "john",
  "lastname": "smith",
  "email": "new_test@email.com",
  "groupuuid": "00000000-0000-0000-0000-000000000000",
  "role": "member"
}
```
   
##### Response

***
* Headers
   
```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```
   
* Body
```
{
  "id": "000000000000000000000000",
  "uuid": "00000000-0000-0000-0000-000000000000",
  "username": "newUserName",
  "firstname": "john",
  "lastname": "smith",
  "email": "new_test@email.com",
  "role": "member",
  "groupuuid": "00000000-0000-0000-0000-000000000000",
  "last_modified": "2019-07-19 20:22:44.322303274 +0000 UTC",
  "creation_datetime": "2019-06-07 20:22:44.322303274 +0000 UTC"
}
```
   
   
#### 4. Delete User
* DELETE - /users/{userId}
   
##### Request

***
* Headers
   
```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```
   
##### Response

***
* Headers
   
```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```
   
### IV) User Group Routes (Admins Only)

___
#### 1. List User Group(s)
* GET - /groups/{groupId}
* groupId parameter is optional, if used the request will only return an object for that item.
  
##### Request

***
* Headers
  
```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```
   
##### Response

***
* Headers
   
```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```
   
* Body
```
[
  {
    "id": "000000000000000000000000",
    "grouptype": "normal",
    "uuid": "00000000-0000-0000-0000-000000000000",
    "name": "groupName",        
    "last_modified": "2019-06-07 20:17:14.358617998 +0000 UTC",
    "creation_datetime": "2019-06-07 20:17:14.358617998 +0000 UTC"
  }
]
```
   
   
#### 2. Create User Group
* POST - /groups
   
##### Request

***
* Headers
   
```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```
   
* Body
```
{
  "Name": "newGroup"
}
```
   
##### Response

***
* Headers
   
```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```
   
* Body
```
{
  "id": "000000000000000000000000",
  "grouptype": "normal",
  "uuid": "00000000-0000-0000-0000-000000000000",
  "name": "newGroup",    
  "last_modified": "2019-06-07 20:18:15.145971952 +0000 UTC",
  "creation_datetime": "2019-06-07 20:18:15.145971952 +0000 UTC"
}
```

   
#### 3. Modify User Group
* PATCH - /groups/{groupId}
   
##### Request

***
* Headers
   
```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```
   
* Body
```
{
  "Name": "newGroupName"
}
```
   
##### Response

***
* Headers
   
```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```
   
* Body
```
{
  "id": "000000000000000000000000",
  "grouptype": "normal",
  "uuid": "00000000-0000-0000-0000-000000000000",
  "name": "newGroupName",    
  "last_modified": "2019-06-07 20:18:15.145971952 +0000 UTC",
  "creation_datetime": "2019-06-07 20:18:15.145971952 +0000 UTC"
}
```
   
#### 4. Delete User Group
* DELETE - /groups/{groupId}
   
##### Request

***
* Headers
   
```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```
   
##### Response

***
* Headers
   
```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```
