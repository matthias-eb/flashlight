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

	fmt.Printf("Found User in Database: %+v\n", u)
	uJSON, err := json.Marshal(u[0])
	if err != nil {
		return user, err
	}
	err = json.Unmarshal(uJSON, &user)
	fmt.Printf("After Marshal and Unmarshal: %+v\n", user)
	if err != nil {
		fmt.Printf("Error while Unmarshaling: %+v", err.Error())
		return user, err
	}
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
