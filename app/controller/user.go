package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	mw "github.com/matthias-eb/flashlight/app/middleware"
	db "github.com/matthias-eb/flashlight/app/model"
	st "github.com/matthias-eb/flashlight/app/structs"
)

//Login is for logging in a User. It can receive a POST or a GET Method.
func Login(w http.ResponseWriter, r *http.Request) {
	mw.SetupSession(w, r)

	fmt.Println("Request for /login coming in: " + r.Method)

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := db.GetUser(username)

		if err == nil {
			fmt.Printf("Checking %+v with %+v\n", password, user.Password)
			err = mw.AuthenticateUser(username, password, user.Password)
			if err == nil {
				mw.SaveSession(w, r)
			}
		}

		if err != nil {
			mw.Templ.ExecuteTemplate(w, "login.tmpl", st.Data{Title: "Flashlight Login", Error: []string{"Benutzername oder Password waren falsch"}})
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else if r.Method == "GET" {
		mw.Templ.ExecuteTemplate(w, "login.tmpl", st.Data{Title: "Flashlight Login", Error: nil})
	}
}

//Logout logs out a User.
func Logout(w http.ResponseWriter, r *http.Request) {
	mw.SetupSession(w, r)
	fmt.Println("Starting logout")
	mw.EndSession(w, r)
	mw.SaveSession(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//Register registrates a User if the username isn't taken and the password isn't too short and it matches the password confirmation.
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		passwordConfirm := r.FormValue("password_confirm")

		var errors []string

		if password != passwordConfirm {
			errors = append(errors, "Passwörter stimmen nicht überein")
		}

		if len(password) < 8 {
			errors = append(errors, "Die Länge des Passworts muss mindestens 8 Zeichen sein")
		}

		//Potentially dangerous due to SQL injection, gotta check!
		if _, err := db.GetUser(username); err == nil {
			errors = append(errors, "Benutzername existiert bereits")
		}
		mw.SetupSession(w, r)

		if len(errors) > 0 {
			mw.SaveSession(w, r)
			mw.Templ.ExecuteTemplate(w, "register.tmpl", st.Data{Title: "Flashlight Registrieren", Error: errors})
			return
		}

		passwordHashed, err := mw.HashPassword(password)
		if err != nil {
			fmt.Println(err)
			errors = append(errors, "Passwort konnte nicht gehasht werden.")
		}

		if len(errors) > 0 {
			mw.Templ.ExecuteTemplate(w, "register.tmpl", st.Data{Title: "Flashlight Registration", Error: errors})
			return
		}

		user := db.User{
			Name:     username,
			Type:     "user",
			Password: passwordHashed,
		}

		err = db.AddUser(user)
		if err != nil {
			fmt.Printf("Error while Adding User: %+v\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		user, err = db.GetUser(username)
		if err != nil {
			fmt.Printf("Error Getting User: %+v\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		mw.AuthenticateUser(username, password, user.Password)
		mw.SaveSession(w, r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if r.Method == "GET" {
		mw.Templ.ExecuteTemplate(w, "register.tmpl", st.Data{Title: "Flashlight Registrierung", Error: nil})
	}

}

// UploadImage checks the User for a valid Session, then saves up to 100MB in filesize to the filesystem and saves everything necessary to the Database
func UploadImage(w http.ResponseWriter, r *http.Request) {
	var data st.Data      //Data to be added to the template
	var errorstr []string //Any Error output meant for the user gets saved here

	//Check if User is logged in
	mw.SetupSession(w, r)
	username, err := mw.CheckAuthentication(w, r)
	if err != nil {
		fmt.Printf("Error while authenticating: %+v\n", err.Error())
		mw.EndSession(w, r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data = st.Data{
		Title: "Flashlight Image Upload",
		User:  username,
	}

	if r.Method == "GET" {

		err = mw.Templ.ExecuteTemplate(w, "add_image.tmpl", data)
		if err != nil {
			fmt.Printf("Error while creating Template: %+v\n", err.Error())
			fmt.Println(err)
		}
		return
	}

	// Parse up to 100 Megabytes of filesize
	err = r.ParseMultipartForm(100000000)
	if err != nil {
		fmt.Printf("Error in image Upload: %+v\n", err.Error())
		errorstr = append(errorstr, "Image too large. Please keep it under 100 MB")
	}
	file, handler, err := r.FormFile("newImage")
	if err != nil {
		fmt.Printf("Error in image Upload: %+v\n", err.Error())
		errorstr = append(errorstr, "Image File not Found in Request")
		data.Error = errorstr
		mw.Templ.ExecuteTemplate(w, "add_image.tmpl", data)
		return
	}

	defer file.Close()

	//Create a Temporary File to get a random File name.
	imageFile, err := ioutil.TempFile(os.TempDir(), "upload-*.png")
	if err != nil {
		fmt.Printf("Creation not possible: %+v\n", err.Error())
		errorstr = append(errorstr, "Error creating File")
	}
	//Use the random Filename for the actual File if it doesn't exist yet
	imageName := imageFile.Name()
	actualImagePath := "images/" + filepath.Base(imageName)
	for _, err := os.Stat(actualImagePath); err == nil; {
		fmt.Printf("This File already exists!")
		imageFile.Close()
		imageFile, err = ioutil.TempFile(os.TempDir(), "upload-*.png")
		if err != nil {
			fmt.Printf("Creation not possible: %+v\n", err.Error())
			errorstr = append(errorstr, "Error creating File")
		}
		imageName = imageFile.Name()
		actualImagePath = "images/" + filepath.Base(imageName)
	}
	fmt.Printf("Writing file to %+v\n", imageName)
	imageFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading Bytes: %+v\n", err.Error())
		errorstr = append(errorstr, "Image not readable")
	}
	err = ioutil.WriteFile(actualImagePath, fileBytes, 0666)

	description := r.FormValue("description")

	fmt.Printf("Upload-Filename: %+v\nDescription: %+v\nSize: %+v\n", actualImagePath, description, handler.Size)

	err = db.AddImage(username, actualImagePath, description)
	if err != nil {
		fmt.Printf("Error while Saving File to Database: %+v\n", err.Error())
		errorstr = append(errorstr, "Uploading information to Database went wrong.")
	}

	if len(errorstr) > 0 {
		data.Error = errorstr
		mw.Templ.ExecuteTemplate(w, "add_image.tmpl", data)
		return
	}

	http.Redirect(w, r, "/my-images", http.StatusSeeOther)
}

//GetImages is the Handler for the User "My Images" Webpage
func GetImages(w http.ResponseWriter, r *http.Request) {
	mw.SetupSession(w, r)
	username, err := mw.CheckAuthentication(w, r)
	if err != nil {
		fmt.Printf("User not logged in.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userImages, err := db.GetImagesForUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := st.Data{
		Images: userImages,
		Title:  "Flashlight Meine Bilder",
		Error:  nil,
		User:   username,
	}

	mw.Templ.ExecuteTemplate(w, "images.tmpl", data)
}

// AddComment saves a Comment that was added to an image
func AddComment(w http.ResponseWriter, r *http.Request) {
	mw.SetupSession(w, r)
	username, err := mw.CheckAuthentication(w, r)
	if err != nil {
		fmt.Printf("User not authenticated while trying to Post comment. Redirectig..\n")
		mw.SaveSession(w, r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	comment := r.FormValue("comment")
	path := r.FormValue("imagepath")

	err = db.AddComment(username, comment, path)
	if err != nil {
		fmt.Println("Error adding Comment:")
		fmt.Println(err)
		http.Error(w, "Error Posting Comment", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/#"+path, http.StatusSeeOther)
}

//LikeImage likes an Image if the User is logged in and has a valid session and the User did not like nor is the owner of the image.
func LikeImage(w http.ResponseWriter, r *http.Request) {
	mw.SetupSession(w, r)
	username, err := mw.CheckAuthentication(w, r)
	if err != nil {
		fmt.Printf("User not authenticated while trying to Post. Redirectig..\n")
		mw.SaveSession(w, r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	path := r.FormValue("imagepath")

	err = db.AddLike(username, path)
	if err != nil {
		fmt.Println("Error while liking Image: ")
		fmt.Println(err)
	}
	mw.SaveSession(w, r)
	http.Redirect(w, r, "/#"+path, http.StatusSeeOther)
}

func DeleteImage(w http.ResponseWriter, r *http.Request) {
	mw.SetupSession(w, r)
	_, err := mw.CheckAuthentication(w, r)
	if err != nil {
		fmt.Println("Tried to delete Image while not being logged in. Redirecting.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	imagepath := r.FormValue("imagepath")

	err = db.DeleteImage(imagepath)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not delete Image for Path "+imagepath, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/my-images", http.StatusSeeOther)
}
