package helper

import (
	"time"

	"github.com/araddon/dateparse"
)

/*
	ConvertTimeToZone method is to convert utc to specified timezone

	@param timeFromGitlab
	@param timezone
*/
func ConvertTimeToZone(timeFromGitlab string, strZone string) (string, error) {
	// Parse an unknown date format, detect the layout.
	dateTime, err := dateparse.ParseAny(timeFromGitlab)
	if err != nil {
		return timeFromGitlab, err
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
