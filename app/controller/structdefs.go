package controller

//Comment is the Comment that gets Inserted by the user
type Comment struct {
	Commentor string
	Comment   string
}
type Image struct {
	Owner       string
	Date        string
	Path        string
	Likes       string
	Liked       bool
	Description string
	Comments    []Comment
}
type Data struct {
	Title  string
	Error  []string
	User   string
	Images []Image
}
