# Go RESTful API Boilerplate

A RESTful API Boilerplate written in Go.

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
* User setup with NOPASSWD: ALL sudo powers

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

```bash
$ cp conf.json.example conf.json
$ vi conf.json
```
___
## Running the API

To start the API run the following:

```bash
$ sh start.sh
```

To stop the API:

```bash
$ sh stop.sh
```
___
## API Route Guide
### I) Authentication Routes

___
#### 1. Signin
* POST - /auth

##### Request

***
* Header

```
{
  Content-Type: application/json
}
```

##### Response

***
* Header

```
{
  Content-Type: application/json; charset=UTF-8,
  Auth-Token: "",
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0
}
```

* Body
```
{
  "id": "000000000000000000000000",
  "uuid": "00000000-0000-0000-0000-000000000000",
  "username": "username",
  "email": "user@example.com",
  "role": "member",
  "groupuuid": "00000000-0000-0000-0000-000000000000",
  "last_modified": "2019-06-07 20:17:14.630917778 +0000 UTC",
  "creation_datetime": "2019-06-07 20:17:14.630917778 +0000 UTC"
}
```

#### 2. Refresh Token
* GET - /auth

##### Request

***
* Header

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

##### Response

***
* Header

```
{
  Content-Type: application/json; charset=UTF-8,
  Auth-Token: "",
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0
}
```

#### 3. Signout
* DELETE - /auth

##### Request

***
* Header

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

##### Response

***
* Header

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0
}
```

#### 4. API Key - Expires after 6 Months
* GET - /auth/api-key

##### Request

***
* Header

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

##### Response

***
* Header

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Auth-Token: "",
  API-Key: ""
}
```

#### 5. Update Password
* POST - /auth/password

##### Request

***
* Header

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
* Header

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0
}
```

### II) Todo Routes

___
#### 1. List Todo(s)
* GET - /todos/{todoId}
* todoId parameter is optional, if used the request will only return an object for that item.

##### Request

***
* Header
```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

##### Response

***
* Header
```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0
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
* Header

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
* Header

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0
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
* Header

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
* Header

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0
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
* Header

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

##### Response

***
* Header

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0
}
```

### III) Users Routes (Admins Only)

___
#### 1. List User(s)
* GET - /users/{userId}
* userId parameter is optional, if used the request will only return an object for that item.
   
##### Request

***
* Header
   
```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```
   
##### Response

***
* Header
   
```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0
}
```
   
* Body
```
[
  {
    "id": "000000000000000000000000",
    "uuid": "00000000-0000-0000-0000-000000000000",
    "username": "userName",
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
* Header
   
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
  "email": "test@email.com",
  "groupuuid": "00000000-0000-0000-0000-000000000000",
  "role": "member"
}
```
   
##### Response

***
* Header
   
```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0
}
```
   
* Body
```
{
  "id": "000000000000000000000000",
  "uuid": "00000000-0000-0000-0000-000000000000",
  "username": "userName",
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
* Header
   
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
  "email": "new_test@email.com",
  "groupuuid": "00000000-0000-0000-0000-000000000000",
  "role": "member"
}
```
   
##### Response

***
* Header
   
```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0
}
```
   
* Body
```
{
  "id": "000000000000000000000000",
  "uuid": "00000000-0000-0000-0000-000000000000",
  "username": "newUserName",
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
* Header
   
```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```
   
##### Response

***
* Header
   
```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0
}
```
   
### IV) User Group Routes (Admins Only)

___
#### 1. List User Group(s)
* GET - /groups/{groupId}
* groupId parameter is optional, if used the request will only return an object for that item.
  
##### Request

***
* Header
  
```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```
   
##### Response

***
* Header
   
```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0
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
* Header
   
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
* Header
   
```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0
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
* Header
   
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
* Header
   
```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0
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
* Header
   
```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```
   
##### Response

***
* Header
   
```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0
}
```