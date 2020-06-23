package controller

type Comment struct {
	User    string
	Comment string
}
type Image struct {
	User        string
	Date        string
	Path        string
	Likes       string
	Description string
	Comments    []Comment
}
type Data struct {
	Title  string
	Error  []string
	User   string
	Images []Image
}
