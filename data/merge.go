package data

/*
	Stuctures of a merge
*/
type Merge struct {
	ObjectKind       string           `json:"object_kind"`
	ObjectAttributes ObjectAttributes `json:"object_attributes"`
	Changes          Changes          `json:"changes"`
}

type ObjectAttributes struct {
	ID              float64 `json:"id"`
	TargetBranch    string  `json:"target_branch"`
	SourceBranch    string  `json:"source_branch"`
	SourceProjectID float64 `json:"source_project_id"`
	AuthorID        float64 `json:"author_id"`
	AssigneeID      float64 `json:"assignee_id"`
	Title           string  `json:"title"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	URL             string  `json:"url"`
	StCommits       float64 `json:"st_commits"`
	StDiffs         float64 `json:"st_diffs"`
	MilestoneID     float64 `json:"milestone_id"`
	State           string  `json:"state"`
	MergeStatus     string  `json:"merge_status"`
	TargetProjectID float64 `json:"target_project_id"`
	Iid             float64 `json:"iid"`
	Description     string  `json:"description"`
	Source          Branch  `json:"source"`
	Target          Branch  `json:"target"`
	LastCommit      Commit  `json:"last_commit"`
}

type Branch struct {
	Name            string  `json:"name"`
	SSHURL          string  `json:"ssh_url"`
	HTTPURL         string  `json:"http_url"`
	VisibilityLevel float64 `json:"visibility_level"`
	Namespace       string  `json:"namespace"`
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
