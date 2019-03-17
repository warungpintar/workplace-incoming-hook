package data

/*
	Stuctures of a merge
*/
type Merge struct {
	Object_kind       string
	Object_attributes Object_attributes
	Changes           Changes `json:"changes"`
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

// Changes object shape from Gitlab payload
type Changes struct {
	Labels Labels `json:"labels"`
}

// Labels object shape from Gitlab payload
type Labels struct {
	Previous []Label `json:"previous"`
	Current  []Label `json:"current"`
}

// Label object shape from Gitlab payload
type Label struct {
	ID          int64       `json:"id"`
	Title       string      `json:"title"`
	Color       string      `json:"color"`
	ProjectID   interface{} `json:"project_id"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	Template    bool        `json:"template"`
	Description string      `json:"description"`
	Type        string      `json:"type"`
	GroupID     int64       `json:"group_id"`
}
