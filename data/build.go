package data

/*
	Stuctures of a build
*/
type Build struct {
	BuildID         float64  `json:"build_id"`
	BuildStatus     string   `json:"build_status"`
	BuildStartedAt  string   `json:"build_started_at"`
	BuildFinishedAt string   `json:"build_finished_at"`
	ProjectID       float64  `json:"project_id"`
	ProjectName     string   `json:"project_name"`
	GitlabURL       string   `json:"gitlab_url"`
	Ref             string   `json:"ref"`
	Sha             string   `json:"sha"`
	BeforeSha       string   `json:"before_sha"`
	PushData        PushData `json:"push_data"`
}

type PushData struct {
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
