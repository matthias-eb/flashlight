package model

import (
	"encoding/json"
	"errors"
	"fmt"
)

//User is a Database conform User
type User struct {
	ID       string `json:"_id"`
	Rev      string `json:"_rev"`
	Type     string `json:"type"`
	Name     string `json:"username"`
	Password string `json:"password"`
}

//GetUser searches in the couchdb for the User with the given username and returns a Map ready to be applied to a Html site
func GetUser(username string) (user User, err error) {
	user = User{}

	query := `
	{
		"selector": {
			"type": "user",
			"username": "%s"
		}
	}`

	u, err := couchDB.QueryJSON(fmt.Sprintf(query, username))
	if err != nil {
		return
	} else if len(u) == 0 {
		err = errors.New("Username not found")
		return
	}

	uJSON, err := json.Marshal(u)
	if err != nil {
		return user, err
	}
	json.Unmarshal(uJSON, &user)

	return user, nil
}

//AddUser searches in the couchdb for the User with the given username and returns a Map ready to be applied to a Html site
func AddUser(user User) (err error) {

	var userMap map[string]interface{}
	uJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	json.Unmarshal(uJSON, &userMap)

	delete(userMap, "_id")
	delete(userMap, "_rev")
	_, _, err = couchDB.Save(userMap, nil)
	if err != nil {
		return err
	}

	return nil
}
