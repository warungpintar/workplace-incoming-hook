package data

/*
	Stuctures of a push
*/
type Push struct {
	Before            string     `json:"before"`
	After             string     `json:"after"`
	Ref               string     `json:"ref"`
	UserID            float64    `json:"user_id"`
	UserName          string     `json:"user_name"`
	ProjectID         float64    `json:"project_id"`
	Repository        Repository `json:"repository"`
	Commits           []Commit   `json:"commits"`
	TotalCommitsCount float64    `json:"total_commits_count"`
}
