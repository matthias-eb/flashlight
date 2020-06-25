package model

import (
	"encoding/json"
	"errors"
	"fmt"

	st "github.com/matthias-eb/flashlight/app/structs"
)

//User is a Database conform User
type User struct {
	ID       string `json:"_id"`
	Rev      string `json:"_rev"`
	Type     string `json:"type"`
	Name     string `json:"username"`
	Password string `json:"password"`
}
type image struct {
	ID          string `json:"_id"`
	Rev         string `json:"_rev"`
	Type        string `json:"type"`
	Path        string `json:"path"`
	UserID      string `json:"user"`
	Description string `json:"comment"`
	Date        string `json:"timestamp"`
}
type comment struct {
	ID        string `json:"_id"`
	Rev       string `json:"_rev"`
	Type      string `json:"type"`
	UserID    string `json:"user"`
	Comment   string `json:"comment"`
	Timestamp string `json:"timestamp"`
	ImageID   string `json:"parent"`
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

//AddImage Saves an Image to the Database and returns an eror if it it wasn't successful
func AddImage(username string, filepath string, description string) (err error) {
	var imageMap map[string]interface{}
	var user User

	//Query Database for User so we can get his ID
	query := `
	{
		"selector": {
			"type": "user",
			"username": "%s"
		}
	}`

	//More than one User and no User found will give an error
	u, err := couchDB.QueryJSON(fmt.Sprintf(query, username))
	if err != nil {
		return err
	} else if len(u) == 0 {
		err = errors.New("Username not found")
		return
	} else if len(u) > 1 {
		err = errors.New("More than one User with the same Name")
	}

	//Move all Data from the Database into the User object
	uJSON, err := json.Marshal(u[0])
	if err != nil {
		return err
	}
	err = json.Unmarshal(uJSON, &user)
	if err != nil {
		return err
	}

	//Fill image Object with Data
	image := image{
		Type:        "image",
		Path:        filepath,
		Description: description,
		UserID:      user.ID,
	}
	//Create a json map from it
	uJSON, err = json.Marshal(image)
	if err != nil {
		return err
	}
	err = json.Unmarshal(uJSON, &imageMap)
	if err != nil {
		return err
	}
	delete(imageMap, "_id")
	delete(imageMap, "_rev")

	//Save it to the Database
	_, _, err = couchDB.Save(imageMap, nil)
	if err != nil {
		return err
	}
	return nil
}

//GetImagesForUser returns all Images that belong to a user
func GetImagesForUser(username string) (images []st.Image, err error) {
	var imageMap []map[string]interface{}
	var imagesDB []image
	user, err := getUserRaw(username)
	if err != nil {
		fmt.Printf("Error getting User: %+v\n", err)
		return
	}

	query := `
	{
		"selector": {
			"type": "image",
			"user": "%s"
		}
	}`

	imageMap, err = couchDB.QueryJSON(fmt.Sprintf(query, user.ID))
	if err != nil {
		return
	}
	imBytes, err := json.Marshal(imageMap)
	if err != nil {
		return
	}
	err = json.Unmarshal(imBytes, &imagesDB)
	if err != nil {
		return nil, err
	}

	var image st.Image
	images = make([]st.Image, len(imagesDB))
	//Fill return values
	for i, imageDB := range imagesDB {
		image = st.Image{
			Owner:       username,
			Path:        imageDB.Path,
			Description: imageDB.Description,
			Date:        imageDB.Date,
			Liked:       true,
		}

		commentsDB, err := getCommentsForImage(imageDB.ID)
		comments := make([]st.Comment, len(commentsDB))
		if err != nil {
			return nil, err
		}
		for j, commentDB := range commentsDB {
			comment := st.Comment{
				Comment:   commentDB.Comment,
				Timestamp: commentDB.Timestamp,
			}
			commentor, err := getUserFromID(commentDB.UserID)
			if err != nil {
				return nil, err
			}
			comment.Commentor = commentor.Name
			comments[j] = comment
		}

		image.Comments = comments
		//image.NrComments = strconv.Itoa(len(comments))
		image.NrComments = len(comments)
		fmt.Printf("Nr of Comments: %+v\n", image.NrComments)

		//ToDo: Likes
		image.Likes = 12

		images[i] = image
	}
	return
}

func getUserRaw(username string) (user User, err error) {
	var u []map[string]interface{}
	var uJSON []byte
	query := `
	{
		"selector": {
			"type": "user",
			"username": "%s"
		}
	}`

	u, err = couchDB.QueryJSON(fmt.Sprintf(query, username))
	if err != nil {
		return
	} else if len(u) == 0 {
		err = errors.New("Username not found")
		return
	} else if len(u) > 1 {
		err = errors.New("More than one User")
		return
	}

	uJSON, err = json.Marshal(u[0])
	if err != nil {
		return
	}
	err = json.Unmarshal(uJSON, &user)
	return
}

func getCommentsForImage(imageID string) (comments []comment, err error) {
	query := `
	{
		"selector": {
			"type": "comment",
			"parent": "%s"
		}
	}`

	commentsMap, err := couchDB.QueryJSON(fmt.Sprintf(query, imageID))
	if err != nil {
		return
	}
	cJSON, err := json.Marshal(commentsMap)
	if err != nil {
		return
	}
	json.Unmarshal(cJSON, &comments)
	return
}

func getUserFromID(userID string) (user User, err error) {
	query := `
	{
		"selector": {
			"type": "user",
			"_id": "%s"
		}
	}`

	userMap, err := couchDB.QueryJSON(fmt.Sprintf(query, userID))
	if err != nil {
		return
	}
	cJSON, err := json.Marshal(userMap)
	if err != nil {
		return
	}
	json.Unmarshal(cJSON, &user)
	return
}
