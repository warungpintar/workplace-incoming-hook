package data

type Repository struct {
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
}

type Commit struct {
	Id        string `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Url       string `json:"url"`
	Author    Author `json:"author"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
