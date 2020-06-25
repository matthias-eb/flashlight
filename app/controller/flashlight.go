package controller

import (
	"fmt"
	"net/http"

	mw "github.com/matthias-eb/flashlight/app/middleware"
	st "github.com/matthias-eb/flashlight/app/structs"
)

//Preview responds to a Get Request to Root. It will then show the index Page with the newest Posts of all Users
func Preview(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Generating Preview")
	mw.SetupSession(w, r)
	username, err := mw.CheckAuthentication(w, r)
	data := st.Data{
		Title: "Flashlight",
		Error: nil,
	}
	var images []st.Image
	var comments []st.Comment
	comment := st.Comment{
		Commentor: "Alex",
		Comment:   "This is a dang stupid comment, but haven't I seen that picture before in my life?",
	}
	comments = append(comments, comment)

	image := st.Image{
		Owner:       "Max Mustermann",
		Date:        "23.10.2017 - 15:00",
		Path:        "images/new-york-taxi.jpg",
		Likes:       10,
		Description: "Some quick example",
		Comments:    comments,
	}
	images = append(images, image)

	comment = st.Comment{
		Commentor: "Alex",
		Comment:   "This is a dang stupid comment, but haven't I seen that picture before in my life?",
	}
	comments = append(comments, comment)
	comment = st.Comment{
		Commentor: "Ben",
		Comment:   "This is a dang stupid comment, but haven't I seen that comment before in my life?",
	}
	comments = append(comments, comment)

	image = st.Image{
		Owner:       "Max Mustermann",
		Date:        "23.10.2017 - 14:00",
		Path:        "images/new-york-taxi.jpg",
		Likes:       10,
		Description: "Some quick example",
		Comments:    comments,
	}
	images = append(images, image)
	if err != nil {
		data.User = ""
		fmt.Printf("User is not logged in.")
	} else {
		fmt.Printf("User %+v is logged in.\n", username)
		data.User = username
		images[0].Liked = true
	}
	data.Images = images
	err = mw.Templ.ExecuteTemplate(w, "index.tmpl", data)
	if err != nil {
		fmt.Printf("Error executing Templates: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
