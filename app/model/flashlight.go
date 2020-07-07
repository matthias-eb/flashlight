package model

import (
	"encoding/json"
	"time"

	st "github.com/matthias-eb/flashlight/app/structs"
)

//GetAllImages returns an array of Images. If the username matches a Like in the Image, Liked is set to true.
//The Images are sorted ascending by timestamp.
func GetAllImages(username string) (images []st.Image, err error) {
	var imagesMap []map[string]interface{}
	var imagesDB []image
	var liked = false
	var currentUser User
	userKnown := username != ""

	query := `{
	"selector": {
		"type": "image"
	},
	"sort": [{
		"timestamp": "desc"
	}]
	}`

	imagesMap, err = couchDB.QueryJSON(query)
	if err != nil {
		return
	}

	imJSON, err := json.Marshal(imagesMap)
	if err != nil {
		return
	}
	err = json.Unmarshal(imJSON, &imagesDB)
	if err != nil {
		return
	}

	if userKnown {
		currentUser, err = getDBUser(username)
		if err != nil {
			return
		}
	}

	for _, im := range imagesDB {
		liked = false

		imgts, err := time.Parse("2006-01-02 15:04:05", im.Date)

		imageOwner, err := getUserFromID(im.UserID)
		if err != nil {
			return nil, err
		}

		if userKnown {
			if contains(im.Likes, currentUser.ID) {
				liked = true
			}
			if currentUser.ID == im.UserID {
				liked = true
			}
		}

		commentsDB, err := getCommentsForImage(im.ID)
		var comments []st.Comment
		for _, cm := range commentsDB {
			u, err := getUserFromID(cm.UserID)
			if err != nil {
				return nil, err
			}

			cmts, err := time.Parse("2006-01-02 15:04:05", cm.Timestamp)

			comment := st.Comment{
				Commentor: u.Name,
				Comment:   cm.Comment,
				Timestamp: cmts.Format("2.1.2006 15:04") + " Uhr", //ToDO
			}
			comments = append(comments, comment)
		}

		imageGrenz := st.Image{
			Owner:       imageOwner.Name,
			Date:        imgts.Format("2.1.2006 15:04") + " Uhr",
			Description: im.Description,
			Path:        im.Path,
			Comments:    comments,
			NrComments:  len(comments),
			Likes:       len(im.Likes),
			Liked:       liked,
		}
		images = append(images, imageGrenz)
	}
	return
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
