package main

import (
	"fmt"

	"github.com/warungpintar/workplace-incoming-hook/data"
	"github.com/warungpintar/workplace-incoming-hook/helper"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/nurza/logo"

	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/*
	Global variables
*/
var (
	// Logging
	l logo.Logger

	// Configuration
	ThreadGitlab        string
	ThreadTuleap        string
	ThreadAppCenter     string
	PushIcon            string // Push icon (Fb emoji)
	MergeIcon           string // Merge icon (Fb emoji)
	BuildIcon           string // Build icon (Fb emoji)
	BotStartMessage     string // Bot's start message
	FbAPIUrl            string // Fb API URL
	Verbose             bool   // Enable verbose mode
	ShowAllCommits      bool   // Show all commits rather than latest
	HTTPTimeout         int    // Http timeout in second
	ChatType            string
	TuleapURL           string
	Port                string
	URLNoteHookFunction string
	TimeZone            string

	// Misc
	currentBuildID float64   // Current build ID
	n              = "%5CnX" // Encoded line return
)

type GitlabServ struct{}

/*
	Flags
*/
var (
	ConfigFile = flag.String("f", "config.json", "Configuration file")
)

/*
	Load configuration file
*/
func LoadConf() {
	conf := struct {
		ThreadGitlab        string
		ThreadAppCenter     string
		ThreadTuleap        string
		PushIcon            string
		MergeIcon           string
		BuildIcon           string
		BotStartMessage     string
		FbAPIUrl            string
		Verbose             bool
		ShowAllCommits      bool
		HTTPTimeout         float64
		ChatType            string
		TuleapURL           string
		Port                string
		URLNoteHookFunction string
		TimeZone            string
	}{}

	content, err := ioutil.ReadFile(*ConfigFile)
	if err != nil {
		l.Critical("Error: Read config file error: " + err.Error())
	}

	err = json.Unmarshal(content, &conf)
	if err != nil {
		l.Critical("Error: Parse config file error: " + err.Error())
	}

	PushIcon = conf.PushIcon
	MergeIcon = conf.MergeIcon
	BuildIcon = conf.BuildIcon
	BotStartMessage = conf.BotStartMessage
	FbAPIUrl = conf.FbAPIUrl
	Verbose = conf.Verbose
	ShowAllCommits = conf.ShowAllCommits
	HTTPTimeout = int(conf.HTTPTimeout)
	ChatType = conf.ChatType
	TuleapURL = conf.TuleapURL
	Port = conf.Port
	ThreadGitlab = conf.ThreadGitlab
	ThreadAppCenter = conf.ThreadAppCenter
	ThreadTuleap = conf.ThreadTuleap
	URLNoteHookFunction = conf.URLNoteHookFunction
	TimeZone = conf.TimeZone
}

/*
	HTTP POST request

	target:		url target
	payload:	payload to send

	Returned values:

	int:	HTTP response status code
	string:	HTTP response body
*/
func Post(target string, payload string) (int, string) {
	// Variables
	var err error          // Error catching
	var res *http.Response // HTTP response
	var req *http.Request  // HTTP request
	var body []byte        // Body response

	// Build request
	l.Debug(bytes.NewBufferString(payload))
	req, _ = http.NewRequest("POST", target, bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	// Do request
	client := &http.Client{}
	client.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   time.Duration(HTTPTimeout) * time.Second,
			KeepAlive: time.Duration(HTTPTimeout) * time.Second,
		}).Dial,
		TLSHandshakeTimeout: time.Duration(HTTPTimeout) * time.Second,
	}

	res, err = client.Do(req)

	if err != nil {
		l.Error("Error : Curl POST : " + err.Error())
		if res != nil {
			return res.StatusCode, ""
		}

		return 0, ""
	}
	defer res.Body.Close()

	// Read body
	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		l.Error("Error : Curl POST body read : " + err.Error())
	}

	return res.StatusCode, string(body)
}

/*
	Encode the git commit message with replacing some special characters not allowed by the Fb API

	@param origin Git message to encode
*/

func MessageEncodeX(origin string) string {
	return strings.Replace(origin, "%5CnX", "\\n\\n", -1)
}

func MessageEncode(origin string) string {
	var result string

	for _, e := range strings.Split(origin, "") {
		switch e {
		case "\n":
			result += "%5CnX"
		case "&":
			result += " and "
		default:
			result += e
		}
	}
	return result
}

/*
	Send a message on WorkChat

	@param channel : Targeted channel (could be personal or group)
*/
func SendWorkchatMessage(channel, message string, chattype string) {
	// Variables
	var payload string // POST data sent to Fb
	// var icon string    // Fb emoji

	// toLower(channel)
	l.Silly("toLower =", channel)
	channel = strings.ToLower(channel)
	l.Silly("toLower =", channel)

	// POST Payload formating
	payload = ""
	if chattype == "group" {
		payload += `{"recipient": { "thread_key": "` + strings.ToLower(channel) + `"} , "message": { "text": "` + message + `"}}`
	} else {
		payload += `{"recipient": { "id": "` + strings.ToLower(channel) + `"} , "message": { "text": "` + message + `"}}`
	}

	// Debug information
	if Verbose {
		l.Debug("payload =", payload)
	}

	code, body := Post(FbAPIUrl, MessageEncodeX(payload))
	if code != 200 {
		l.Error("Error post, Fb API returned:", body)
	}

	// Debug information
	if Verbose {
		l.Debug("Fb API returned:", body)
	}
}

/*
	Handler function to handle http requests for push

	@param w http.ResponseWriter
	@param r *http.Request
*/
func (s *GitlabServ) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer // Buffer to get request body
	var body string         // Request body (it's a json)

	// Log
	l.Info("Request")

	// Read http request header
	gitlabEvent := r.Header.Get("X-Gitlab-Event")
	if Verbose {
		l.Debug("Gitlab Event =", gitlabEvent)
	}

	// Read service params
	serviceParam := r.URL.Query().Get("service")
	if Verbose {
		l.Debug("Service =", serviceParam)
	}

	// Read http request body and put it in a string
	if _, err := buffer.ReadFrom(r.Body); err != nil {
		l.Error("Error : Read http request body failed :", err)
	}
	body = buffer.String()

	// Debug information
	if Verbose {
		l.Debug("JsonString receive =", body)
	}

	switch serviceParam {
	case "tuleap":
		TaskHandler(body)
	case "appcenter":
		AppCenterHandler(body)
	default:
		switch gitlabEvent {
		case "Push Hook":
			PushHandler(body)
		case "Merge Request Hook":
			MergeHandler(body)
		case "Build Hook":
			BuildHandler(body)
		case "Note Hook":
			CommentHandler(body)
		}
	}
}

// CommentHandler call cloud function which handle comment event on gitlab
// and send to workplace bot
func CommentHandler(body string) {

	// whether UrlNoteHookFunction has been set #see on config.json
	if URLNoteHookFunction == "" {
		l.Error("url comment service not ")
		return
	}

	var json = []byte(body)
	req, _ := http.NewRequest("POST", URLNoteHookFunction, bytes.NewBuffer(json))
	req.Header.Set("X-Gitlab-Event", "Note Hook")
	req.Header.Set("Content-Type", "application/json")

	l.Info("Call service note hook")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		l.Error("Error : call [POST] microservice failed :", err)
	}
	defer resp.Body.Close()

	l.Info("Response status :", resp.Status)
	l.Info("Response header :", resp.Header)
}

func PushHandler(body string) {
	var j data.Push
	var err error         // Error catching
	var message string    // Bot's message
	var dateString string // Time of the last commit

	// Parse json and put it in a the data.Build structure
	err = json.Unmarshal([]byte(body), &j)
	if err != nil {
		// Error
		l.Error("Error : Json parser failed :", err)
	} else {
		// Ok
		// Debug information
		if Verbose {
			l.Debug("JsonObject =", j)
		}

		// Build the message

		// Date parsing (parsing result example : 16 Jun 19 20:18)
		dateString, _ = helper.ConvertTimeToZone(j.Commits[0].Timestamp, TimeZone)

		// Message
		lastCommit := j.Commits[len(j.Commits)-1]
		commitCount := strconv.FormatFloat(j.TotalCommitsCount, 'f', 0, 64)
		if ShowAllCommits {
			message += "Push on *" + j.Repository.Name + "* by *" + j.UserName + "* at *" + dateString + "* on branch *" + j.Ref + "*:" + n // First line
			message += commitCount + " commits :"                                                                                           // Second line
			for i := range j.Commits {
				c := j.Commits[i]
				message += n + "< " + c.URL + " | " + c.ID[0:7] + " >: " + "_" + MessageEncode(c.Message) + "_"
			}
		} else {
			// First line
			message += "[PUSH] " + n
			message += fmt.Sprintf("Push on *%s* by *%s* at *%s* on branch *%s*: ",
				j.Repository.Name,
				j.UserName,
				dateString,
				j.Ref) + n
			// Second line
			message += "Last commit : < " + lastCommit.URL + " | " + lastCommit.ID + " > :" + n
			// Third line (last commit message)
			message += "```" + MessageEncode(lastCommit.Message) + "```"
		}
		SendWorkchatMessage(ThreadGitlab, message, ChatType)
	}
}

/*
	Handler function to handle http requests for merge

	@param body string
*/
func MergeHandler(body string) {
	var j data.Merge
	var err error         // Error catching
	var message string    // Bot's message
	var dateString string // Time of the last commit

	// Parse json and put it in a the data.Build structure
	err = json.Unmarshal([]byte(body), &j)
	if err != nil {
		// Error
		l.Error("Error : Json parser failed :", err)
	} else {
		// Ok
		// Debug information
		if Verbose {
			l.Debug("JsonObject =", j)
		}

		// Build the message

		// Date parsing (parsing result example : 16 Jun 19 20:18)
		dateString, _ = helper.ConvertTimeToZone(j.ObjectAttributes.CreatedAt, TimeZone)

		// Message
		// First line
		message += fmt.Sprintf("[MERGE REQUEST %s] ", strings.ToUpper(j.ObjectAttributes.State)) + n
		message += fmt.Sprintf("Target : *%s/%s* Source : *%s/%s* at *%s* ",
			j.ObjectAttributes.Target.Name,
			j.ObjectAttributes.TargetBranch,
			j.ObjectAttributes.Source.Name,
			j.ObjectAttributes.SourceBranch,
			dateString) + n

		// Second line (URL link for merge request location)
		message += fmt.Sprintf("Link : *%s*", j.ObjectAttributes.URL) + n

		// Third Line (Description of merge request)
		message += "Description: " + MessageEncode(j.ObjectAttributes.Description)

		if len(j.Changes.Labels.Current) > 0 || len(j.Changes.Labels.Previous) > 0 {
			message += n + " [LABELS] "
			for _, currentLabel := range j.Changes.Labels.Current {
				added := true
				for _, previousLabel := range j.Changes.Labels.Previous {
					if currentLabel.ID == previousLabel.ID {
						added = false
						break
					}
				}

				if added {
					message += n + "`" + currentLabel.Title + "`" + " " + "*Added*"
				}
			}

			for _, previousLabel := range j.Changes.Labels.Previous {
				removed := true
				for _, currentLabel := range j.Changes.Labels.Current {
					if previousLabel.ID == currentLabel.ID {
						removed = false
						break
					}
				}

				if removed {
					message += n + "`" + previousLabel.Title + "`" + " " + "*Removed*"
				}
			}
		}

		SendWorkchatMessage(ThreadGitlab, message, ChatType)
	}
}

/*
	Handler function to handle http requests for build

	@param body string
*/
func BuildHandler(body string) {
	var j data.Build
	var err error      // Error catching
	var message string // Bot's message
	var date time.Time // Time of the last commit

	// Parse json and put it in a the data.Build structure
	err = json.Unmarshal([]byte(body), &j)
	if err != nil {
		// Error
		l.Error("Error : Json parser failed :", err)
	} else {
		// Ok
		// Debug information
		if Verbose {
			l.Debug("JsonObject =", j)
		}

		// Test if the message is already sent
		if currentBuildID < j.BuildID {
			// Not sent
			currentBuildID = j.BuildID // Update current build ID

			// Send the message

			// Date parsing (parsing result example : 18 November 2014 - 14:34)
			date, _ = time.Parse("2006-01-02T15:04:05Z07:00", j.PushData.Commits[0].Timestamp)
			var dateString = strconv.Itoa(date.Day()) + " " + date.Month().String() + " " + strconv.Itoa(date.Year()) +
				" - " + strconv.Itoa(date.Hour()) + ":" + strconv.Itoa(date.Minute())

			// Message
			lastCommit := j.PushData.Commits[len(j.PushData.Commits)-1]
			// First line
			message += "[BUILD] " + n
			message += fmt.Sprintf("%s : Push on *%s* by *%s* at *%s* on branch *%s*:",
				strings.ToUpper(j.BuildStatus),
				j.PushData.Repository.Name,
				j.PushData.UserName,
				dateString,
				j.Ref) + n
			// Second line
			message += "Last commit : <" + lastCommit.URL + "|" + lastCommit.ID + "> :" + n
			// Third line (last commit message)
			message += "```" + MessageEncode(lastCommit.Message) + "```"
			SendWorkchatMessage(ThreadGitlab, message, ChatType)
		}
	}

}

/*
	Handler function to handle http requests for build

	@param body string
*/

func TaskHandler(body string) {
	var j data.Tptask
	var Task data.TuleapTask
	var err error      // Error catching
	var message string // Bot's message
	var date time.Time // Time of the last commit

	// Parse json and put it in a the data.Build structure
	payload := strings.Split(body, "payload=")
	parsedValue, _ := url.QueryUnescape(payload[1])

	err = json.Unmarshal([]byte(parsedValue), &j)
	if err != nil {
		// Error
		l.Error("Error : Json parser failed :", err)
	} else {
		// Ok
		// Debug information
		if Verbose {
			l.Debug("JsonObject =", j)
		}
		Task.Name = j.User.RealName
		for _, val := range j.Current.Values {
			switch val.Label {
			case "Task title":
				Task.TaskTitle = val.Value.(string)
			case "Status":
				Task.Status = val.VValues[0].Label
			case "Links":
				if len(val.ReverseLinks) > 0 {
					Task.ProjectURL = TuleapURL + "projects/" + string(val.ReverseLinks[0].Tracker.Project.ID)
					Task.ProjectName = val.ReverseLinks[0].Tracker.Project.Label
				}
			case "Artifact ID":
				Task.TaskID = strconv.FormatFloat(val.Value.(float64), 'f', 0, 64)
				Task.TrackerURL = TuleapURL + "plugins/tracker/?aid=" + Task.TaskID
			case "Submitted on":
				Task.SubmittedOn = val.Value.(string)
			case "Details":
				Task.Details = strip.StripTags(val.Value.(string))
			case "Type":
				Task.Type = val.VValues[0].Label
			}
		}
		for _, val := range j.Previous.Values {
			if val.Label == "Status" {
				Task.OldStatus = val.VValues[0].Label
			}
		}

		if Task.Status != Task.OldStatus {
			date, _ = time.Parse(time.RFC3339, Task.SubmittedOn)
			var dateString = date.Format("02 Jan 06 15:04")

			// Message
			// First line
			message += fmt.Sprintf("Move Task *#%s* [%s] on Project *%s* by *%s* at *%s* from *%s* to *%s*",
				Task.TaskID,
				Task.TaskTitle,
				Task.ProjectName,
				Task.Name,
				dateString,
				Task.OldStatus,
				Task.Status) + n
			// Third line (last commit message)
			message += "Description: " + MessageEncode(Task.Details) + n
			message += "Task URL : " + Task.TrackerURL
			SendWorkchatMessage(ThreadTuleap, message, ChatType)
		}
	}
}

/*
	Handler function to handle http requests for appcenter

	@param body string
*/
func AppCenterHandler(body string) {
	var j data.AppCenter
	var err error      // Error catching
	var message string // Bot's message
	var date time.Time // Time of the last commit

	// Parse json and put it in a the data.Build structure
	err = json.Unmarshal([]byte(body), &j)
	if err != nil {
		// Error
		l.Error("Error : Json parser failed :", err)
	} else {
		// Ok
		// Debug information
		if Verbose {
			l.Debug("JsonObject =", j)
		}

		// Send the message

		// Date parsing (parsing result example : 18 November 2014 - 14:34)
		date, _ = time.Parse(time.RFC3339, j.SentAt)
		var dateString = date.Format("02 Jan 06 15:04")

		// Message
		message = ""
		if j.DistributionGroupID != "" {
			message += "*Distributed [" + j.AppDisplayName + "]*" + n // First line
			message += "Group ID *" + j.DistributionGroupID + "* on " + dateString + n
			message += "Install Link: " + j.InstallLink // Third line (last commit message)
		} else if j.Reason == "" {
			message += "*" + j.AppName + "* (" + j.OS + ") Branch *" + j.Branch + "*" + n // First line
			message += "Build *#" + j.BuildID + "* [" + j.BuildStatus + "] on " + dateString + n
			message += "URL: " + j.BuildLink // Third line (last commit message)
		} else {
			message += "*Crash!!! [" + j.AppDisplayName + "]*" + n // First line
			message += "Reason *" + j.Reason + "* [" + j.Name + "] on " + dateString + n
			message += "URL: " + j.URL // Third line (last commit message)
		}

		SendWorkchatMessage(ThreadAppCenter, message, ChatType)
	}
}

/*
	Main function
*/
func main() {
	flag.Parse()                                             // Parse flags
	l.AddTransport(logo.Console).AddColor(logo.ConsoleColor) // Configure Logger
	l.EnableAllLevels()                                      // Configure Logger
	LoadConf()                                               // Load configuration
	l.Info(BotStartMessage)                                  // Logging
	l.Error(http.ListenAndServe(":"+Port, &GitlabServ{}))    // Run HTTP server for push hook
}
