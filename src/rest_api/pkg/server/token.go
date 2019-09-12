/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/19/2019
*/

package server

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"rest_api/pkg"
	"rest_api/pkg/configuration"
)

// CreateToken is used to create a new session JWT token
func CreateToken(user root.User, config configuration.Configuration, exp int64) string {
	var MySigningKey = []byte(config.Secret)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["role"] = user.Role
	claims["username"] = user.Username
	claims["uuid"] = user.Uuid
	claims["groupuuid"] = user.GroupUuid
	claims["exp"] = exp
	tokenString, _ := token.SignedString(MySigningKey)
	return tokenString
}

// DecodeJWT is used to decode a JWT token
func DecodeJWT(curToken string, config configuration.Configuration) []string {
	var MySigningKey = []byte(config.Secret)
	// Decode token
	token, err := jwt.Parse(curToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error")
		}
		return []byte(MySigningKey), nil
	})
	if err != nil {
		return []string{"", ""}
	}
	// Determine user based on token
	tokenClaims := token.Claims.(jwt.MapClaims)
	userUuid := tokenClaims["uuid"].(string)
	userRole := tokenClaims["role"].(string)
	groupUuid := tokenClaims["groupuuid"].(string)
	reSlice := []string{userUuid, userRole, groupUuid}
	return reSlice
}
