package data

/*
	Stuctures of a push
*/
type TuleapTask struct {
	Name string
	TaskTitle string
	TaskID string
	Type string
	Status string
	ProjectURL string
	ProjectName string
	TrackerURL string
	SubmittedOn string
	Details string
}

type Tptask struct {
	Action              string `json:"action"`
	User               	User `json:"User"`
	Current				Current `json:"current"`
}

type User struct {
	RealName string `json:"real_name"`
}

type Current struct	{
	Submitted_On string `json:"submitted_on"`
	Values []Values `json:"values"`
}

type ReverseLinks struct {
	ID int `json:"id"`
	Tracker Tracker `json:"tracker"`
}

type Tracker struct {
	ID int `json:"id"`
	Label string `json:"label"`
	Project Project `json:"project"`
}

type Project struct	{
	ID int `json:"id"`
	Label string `json:"label"`
}

type Values struct {
	Label string `json:"label"`
	Reverse_Links []ReverseLinks `json:"reverse_links"`
	VValues []VValues `json:"values"`
	Value interface{} `json:"value"`
}

type VValues struct {
	Label string `json:"label"`
}