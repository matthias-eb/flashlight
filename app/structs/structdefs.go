package structs

//Comment is the Comment that gets Inserted by the user
type Comment struct {
	Commentor string
	Comment   string
	Timestamp string
}

//Image represents an Image with all needed Information
type Image struct {
	Owner       string
	Date        string
	Path        string
	Likes       int
	Liked       bool
	Description string
	Comments    []Comment
	NrComments  int
}

//Data represents all the Data that is needed in any template file called.
type Data struct {
	Title  string
	Error  []string
	User   string
	Images []Image
}
