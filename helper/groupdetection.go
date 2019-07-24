package helper

import (
	"strings"
)

/*
	GroupDetection method is to get group info from url

	@param url

*/
func GroupDetection(url string) string {
	/*
		Split URL base on gitlab.warungpintar.co url
		You can change with your own url

		This format using this url format "git@gitlab.warungpintar.co:back-end/gitlab-test.git"

		Flow :
		1. Split git@gitlab.warungpintar.co:back-end/gitlab-test.git to [ back-end/gitlab-test.git]
		2. Split [ back-end/gitlab-test.git] to [back-end gitlab-test.git]
		3. Get string in slice[0] which is backend

	*/
	link := strings.Split(url, "git@gitlab.warungpintar.co:") //change git@gitlab.warungpintar.co: tou your own url
	link = strings.Split(link[1], "/")
	group := link[0]
	return group
}
