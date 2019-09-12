/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
9/07/2019
*/

package mongodb

import (
	"rest_api/pkg"
)

// BaseModifyUser is a function that setups the base user struct during a user modification request
func BaseModifyUser(user root.User, curUser *userModel) root.User {
	if len(user.Username) == 0 {
		user.Username = curUser.Username
	}
	if len(user.FirstName) == 0 {
		user.FirstName = curUser.FirstName
	}
	if len(user.LastName) == 0 {
		user.LastName = curUser.LastName
	}
	if len(user.Email) == 0 {
		user.Email = curUser.Email
	}
	if len(user.GroupUuid) == 0 {
		user.GroupUuid = curUser.GroupUuid
	}
	if len(user.Role) == 0 {
		user.Role = curUser.Role
	}
	return user
}