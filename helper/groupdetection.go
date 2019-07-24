package helper

import (
	"strings"
)

func GroupDetection(url string) string {
	link := strings.Split(url, "git@gitlab.warungpintar.co:")
	link = strings.Split(link[1], "/")
	group := link[0]
	return group
}
