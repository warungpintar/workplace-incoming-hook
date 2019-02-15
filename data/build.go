package data

/*
	Stuctures of a build
*/
type Build struct {
	Build_id          float64
	Build_status      string
	Build_started_at  string
	Build_finished_at string
	Project_id        float64
	Project_name      string
	Gitlab_url        string
	Ref               string
	Sha               string
	Before_sha        string
	Push_data         Push_Data
}

type Push_Data struct {
	Before              string
	After               string
	Ref                 string
	User_id             float64
	User_name           string
	Project_id          float64
	Repository          Repository
	Commits             []Commit
	Total_commits_count float64
}
