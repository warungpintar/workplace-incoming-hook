package data

type Repository struct {
	Name        string
	Url         string
	Description string
	Homepage    string
}

type Commit struct {
	Id        string
	Message   string
	Timestamp string
	Url       string
	Author    Author
}

type Author struct {
	Name  string
	Email string
}
