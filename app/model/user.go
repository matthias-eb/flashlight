package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

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
	ID          string   `json:"_id"`
	Rev         string   `json:"_rev"`
	Type        string   `json:"type"`
	Path        string   `json:"path"`
	UserID      string   `json:"user"`
	Description string   `json:"comment"`
	Date        string   `json:"timestamp"`
	Likes       []string `json:"likes"`
}
type commentDB struct {
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

//GetImagesForUser returns all Images that belong to a user
func GetImagesForUser(username string) (images []st.Image, err error) {
	var imageMap []map[string]interface{}
	var imagesDB []image
	user, err := getDBUser(username)
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
		ts, err := time.Parse("2006-01-02 15:04:05", imageDB.Date)
		if err != nil {
			return nil, err
		}
		image = st.Image{
			Owner:       username,
			Path:        imageDB.Path,
			Description: imageDB.Description,
			Date:        ts.Format("2.1.2006 15:04") + " Uhr",
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

		//ToDo: Likes
		image.Likes = len(imageDB.Likes)

		images[i] = image
	}
	return
}

func getDBUser(username string) (user User, err error) {
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
	if len(userMap) < 1 {
		err = errors.New("User not found")
		return
	} else if len(userMap) > 1 {
		err = errors.New("User no unique")
		return
	}
	cJSON, err := json.Marshal(userMap[0])
	if err != nil {
		return
	}
	err = json.Unmarshal(cJSON, &user)
	if err != nil {
		return
	}
	return
}

// AddComment uses the imagepath to search for the image corresponding to it and saves a new Comment for the ID of that Image and with the User.
func AddComment(username string, comment string, imagepath string) (err error) {
	var image image
	var imageMap []map[string]interface{}
	var commentStruct commentDB
	var commentMap map[string]interface{}

	query := `
	{
		"selector": {
			"type": "image",
			"path": "%s"
		}
	}`

	imageMap, err = couchDB.QueryJSON(fmt.Sprintf(query, imagepath))
	if err != nil {
		return
	}
	if len(imageMap) < 1 {
		err = errors.New("Image not found for path " + imagepath)
	} else if len(imageMap) > 1 {
		err = errors.New("Too many Images for path " + imagepath)
	}

	// Marshal only the first image
	iJSON, err := json.Marshal(imageMap[0])
	if err != nil {
		return
	}

	err = json.Unmarshal(iJSON, &image)
	if err != nil {
		return
	}

	u, err := GetUser(username)
	if err != nil {
		return
	}

	commentStruct = commentDB{
		Type:      "comment",
		UserID:    u.ID,
		Comment:   comment,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		ImageID:   image.ID,
	}

	cJSON, err := json.Marshal(commentStruct)
	if err != nil {
		return
	}
	json.Unmarshal(cJSON, &commentMap)

	delete(commentMap, "_id")
	delete(commentMap, "_rev")

	_, _, err = couchDB.Save(commentMap, nil)
	if err != nil {
		return
	}
	return nil
}

//AddLike Saves the User that Liked the Image in its Liked array. It will not Like an Image if the Image already contains the ID or if the ID of the user matches the owner of the image.
func AddLike(username string, imagepath string) (err error) {
	var imageMap map[string]interface{}

	u, err := getDBUser(username)
	if err != nil {
		return
	}

	img, err := getImageFromPath(imagepath)
	if err != nil {
		return
	}

	if contains(img.Likes, u.ID) {
		err = errors.New("The User already Liked this image")
		return
	} else if img.UserID == u.ID {
		err = errors.New("The Owner of an Image cannot Like its own image")
		return
	}

	img.Likes = append(img.Likes, u.ID)

	imgData, err := json.Marshal(img)
	if err != nil {
		return
	}
	err = json.Unmarshal(imgData, &imageMap)
	if err != nil {
		return
	}

	_, _, err = couchDB.Save(imageMap, nil)
	if err != nil {
		return
	}
	return
}
