/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/19/2019
*/

package root

// User is a root struct that is used to store the json encoded data for/from a mongodb user doc.
type User struct {
	Id               string `json:"id,omitempty"`
	Uuid             string `json:"uuid,omitempty"`
	Username         string `json:"username,omitempty"`
	Password         string `json:"password,omitempty"`
	Email            string `json:"email,omitempty"`
	Role             string `json:"role,omitempty"`
	GroupUuid        string `json:"groupuuid,omitempty"`
	LastModified     string `json:"last_modified,omitempty"`
	CreationDatetime string `json:"creation_datetime,omitempty"`
}

// UserService is an interface used to manage the relevant user doc controllers
type UserService interface {
	AuthenticateUser(u User) User
	BlacklistAuthToken(authToken string)
	RefreshToken(tokenData []string) User
	UpdatePassword(tokenData []string, CurrentPassword string, newPassword string) User
	UserCreate(u User) User
	UserDelete(id string) User
	UsersFind() []User
	UserFind(id string) User
	UserUpdate(u User) User
	UserDocInsert(u User) User
}
