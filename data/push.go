package data

/*
	Stuctures of a push
*/
type Push struct {
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
