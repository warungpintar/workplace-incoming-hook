package data

/*
	Stuctures of a push
*/
type Push struct {
	Before              string     `json:"before"`
	After               string     `json:"after"`
	Ref                 string     `json:"ref"`
	User_id             float64    `json:"user_id"`
	User_name           string     `json:"user_name"`
	Project_id          float64    `json:"project_id"`
	Repository          Repository `json:"repository"`
	Commits             []Commit   `json:"commits"`
	Total_commits_count float64    `json:"total_commits_count"`
}
