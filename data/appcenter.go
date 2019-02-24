package data

/*
	Stuctures of a AppCenter
*/
type AppCenter struct {
	AppName			string `json:"app_name"`
	Branch			string `json:"branch"`
	BuildStatus		string `json:"build_status"`
	BuildID			string `json:"build_id"`
	BuildLink		string `json:"build_link"`
	BuildReason		string `json:"build_reason"`
	FinishTime		string `json:"finish_time"`
	IconLink		string `json:"icon_link"`
	NSLink			string `json:"notification_settings_link"`
	OS				string `json:"os"`
	StartTime		string `json:"start_time"`
	SourceVersion	string `json:"source_version"`
	SentAt			string `json:"sent_at"`
}
