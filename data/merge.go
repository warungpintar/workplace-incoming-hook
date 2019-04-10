package data

/*
	Stuctures of a merge
*/
type Merge struct {
	Object_kind       string            `json:"object_kind"`
	Object_attributes Object_attributes `json:"object_attributes"`
	Changes           Changes           `json:"changes"`
}

type Object_attributes struct {
	Id                float64 `json:"id"`
	Target_branch     string  `json:"target_branch"`
	Source_branch     string  `json:"source_branch"`
	Source_project_id float64 `json:"source_project_id"`
	Author_id         float64 `json:"author_id"`
	Assignee_id       float64 `json:"assignee_id"`
	Title             string  `json:"title"`
	Created_at        string  `json:"created_at"`
	Updated_at        string  `json:"updated_at"`
	St_commits        float64 `json:"st_commits"`
	St_diffs          float64 `json:"st_diffs"`
	Milestone_id      float64 `json:"milestone_id"`
	State             string  `json:"state"`
	Merge_status      string  `json:"merge_status"`
	Target_project_id float64 `json:"target_project_id"`
	Iid               float64 `json:"iid"`
	Description       string  `json:"description"`
	Source            Branche `json:"source"`
	Target            Branche `json:"target"`
	Last_commit       Commit  `json:"last_commit"`
}

type Branche struct {
	Name             string  `json:"name"`
	Ssh_url          string  `json:"ssh_url"`
	Http_url         string  `json:"http_url"`
	Visibility_level float64 `json:"visibility_level"`
	Namespace        string  `json:"namespace"`
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
