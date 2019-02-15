package data

/*
	Stuctures of a merge
*/
type Merge struct {
	Object_kind       string
	Object_attributes Object_attributes
}

type Object_attributes struct {
	Id                float64
	Target_branch     string
	Source_branch     string
	Source_project_id float64
	Author_id         float64
	Assignee_id       float64
	Title             string
	Created_at        string
	Updated_at        string
	St_commits        float64
	St_diffs          float64
	Milestone_id      float64
	State             string
	Merge_status      string
	Target_project_id float64
	Iid               float64
	Description       string
	Source            Branche
	Target            Branche
	Last_commit       Commit
}

type Branche struct {
	Name             string
	Ssh_url          string
	Http_url         string
	Visibility_level float64
	Namespace        string
}
