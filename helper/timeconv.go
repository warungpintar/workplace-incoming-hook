package helper

import (
	"time"
)

/*
	ConvertTimeToZone method is to convert utc to specified timezone

	@param timeFromGitlab
	@param timezone
*/
func ConvertTimeToZone(timeFromGitlab string, strZone string) (string, error) {

	// adjustment format
	dateTime, err := time.Parse("2006-01-02 15:04:05 UTC", timeFromGitlab)
	if err != nil {
		return timeFromGitlab, err // invalid dateTime format
	}

	// validate time zone
	location, err := time.LoadLocation(strZone)
	if err != nil {
		// if strZone invalid than return back formated origin time with error msg
		return dateTime.Format("02 Jan 06 15:04"), err
	}

	local := dateTime.In(location)
	return local.Format("02 Jan 06 15:04"), nil
}
