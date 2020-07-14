package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

//AddImage Saves an Image to the Database and returns an eror if it it wasn't successful
func AddImage(username string, filepath string, description string) (err error) {
	var imageMap map[string]interface{}

	user, err := getDBUser(username)
	if err != nil {
		return
	}

	ts := time.Now()
	timestamp := ts.Format("2006-01-02 15:04:05")

	//Fill image Object with Data
	image := image{
		Type:        "image",
		Path:        filepath,
		Description: description,
		UserID:      user.ID,
		Date:        timestamp,
	}
	//Create a json map from it
	uJSON, err := json.Marshal(image)
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

func getCommentsForImage(imageID string) (comments []commentDB, err error) {
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
	err = json.Unmarshal(cJSON, &comments)
	if err != nil {
		return
	}
	return
}

func getImageFromPath(imagepath string) (img image, err error) {
	query := `
	{
		"selector": {
			"type": "image",
			"path": "%s"
		}
	}`

	imageMap, err := couchDB.QueryJSON(fmt.Sprintf(query, imagepath))
	if err != nil {
		return
	}
	if len(imageMap) > 1 {
		err = errors.New("Too Many Images for this Path")
		return
	} else if len(imageMap) == 0 {
		err = errors.New("Image not found")
		return
	}
	imBytes, err := json.Marshal(imageMap[0])
	if err != nil {
		return
	}
	err = json.Unmarshal(imBytes, &img)
	if err != nil {
		return
	}
	return
}

//DeleteImage deletes an Image from the Database
func DeleteImage(imagepath string) (err error) {
	var imgComments []commentDB
	img, err := getImageFromPath(imagepath)
	if err != nil {
		return err
	}

	fmt.Println("Deleting Objects for Imagepath " + img.Path + ":")
	fmt.Println("\tImageID: " + img.ID)

	query := `
	{
		"selector": {
			"type": "comment",
			"parent": "%s"
		}
	}
	`

	//Delete all comments belonging to that image.
	commentMap, err := couchDB.QueryJSON(fmt.Sprintf(query, img.ID))
	if err != nil {
		return err
	}
	cBytes, err := json.Marshal(commentMap)
	if err != nil {
		return err
	}

	err = json.Unmarshal(cBytes, &imgComments)
	if err != nil {
		return
	}

	for _, comm := range imgComments {
		fmt.Println("\tComment ID: " + comm.ID)
		err = couchDB.Delete(comm.ID)
		if err != nil {
			return
		}
	}

	err = couchDB.Delete(img.ID)
	return
}
