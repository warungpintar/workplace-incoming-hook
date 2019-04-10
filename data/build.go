package data

/*
	Stuctures of a build
*/
type Build struct {
	Build_id          float64   `json:"build_id"`
	Build_status      string    `json:"build_status"`
	Build_started_at  string    `json:"build_started_at"`
	Build_finished_at string    `json:"build_finished_at"`
	Project_id        float64   `json:"project_id"`
	Project_name      string    `json:"project_name"`
	Gitlab_url        string    `json:"gitlab_url"`
	Ref               string    `json:"ref"`
	Sha               string    `json:"sha"`
	Before_sha        string    `json:"before_sha"`
	Push_data         Push_Data `json:"push_data"`
}

type Push_Data struct {
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
